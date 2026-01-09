package nats

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/byteweap/wukong/plugin/broker"
)

// ID is NATS broker implementation identifier.
const ID = "nats(core)"

// Broker implements broker.Broker with NATS Core.
type Broker struct {
	opts *options
	nc   *nats.Conn
}

var _ broker.Broker = (*Broker)(nil)

// New creates a NATS Core broker and connects immediately.
func New(opts ...Option) (*Broker, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	nc, err := nats.Connect(o.urls, buildNatsOptions(o)...)
	if err != nil {
		return nil, err
	}

	return &Broker{opts: o, nc: nc}, nil
}

// NewWith creates broker with existing *nats.Conn (caller owns its lifecycle).
func NewWith(nc *nats.Conn) *Broker {
	return &Broker{opts: defaultOptions(), nc: nc}
}

func (b *Broker) ID() string { return ID }

func (b *Broker) Publish(ctx context.Context, subject string, data []byte, opts ...broker.PublishOption) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var po broker.PublishOptions
	for _, opt := range opts {
		opt(&po)
	}

	m := &nats.Msg{Subject: subject, Data: data}
	if len(po.Header) > 0 {
		m.Header = toNatsHeader(po.Header)
	}
	if po.Reply != "" {
		m.Reply = po.Reply
	}
	return b.nc.PublishMsg(m)
}

func (b *Broker) Subscribe(ctx context.Context, subject string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscription, error) {
	var so broker.SubscribeOptions
	for _, opt := range opts {
		opt(&so)
	}

	cb := func(m *nats.Msg) {
		msg := &broker.Message{
			Subject: m.Subject,
			Reply:   m.Reply,
			Data:    m.Data,
		}
		if len(m.Header) > 0 {
			msg.Header = fromNatsHeader(m.Header)
		}
		handler(ctx, msg)
	}

	var (
		sub *nats.Subscription
		err error
	)
	if so.Queue != "" {
		sub, err = b.nc.QueueSubscribe(subject, so.Queue, cb)
	} else {
		sub, err = b.nc.Subscribe(subject, cb)
	}
	if err != nil {
		return nil, err
	}

	// ctx 取消时自动退订（避免调用方遗忘）
	if ctx != nil {
		if done := ctx.Done(); done != nil {
			go func() {
				<-done
				_ = sub.Unsubscribe()
			}()
		}
	}

	return &subscription{sub: sub}, nil
}

func (b *Broker) Request(ctx context.Context, subject string, data []byte, opts ...broker.RequestOption) (*broker.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var ro broker.RequestOptions
	for _, opt := range opts {
		opt(&ro)
	}

	req := &nats.Msg{Subject: subject, Data: data}
	if len(ro.Header) > 0 {
		req.Header = toNatsHeader(ro.Header)
	}

	resp, err := b.nc.RequestMsgWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	out := &broker.Message{
		Subject: resp.Subject,
		Reply:   resp.Reply,
		Data:    resp.Data,
	}
	if len(resp.Header) > 0 {
		out.Header = fromNatsHeader(resp.Header)
	}
	return out, nil
}

func (b *Broker) Reply(ctx context.Context, msg *broker.Message, data []byte, opts ...broker.ReplyOption) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if msg == nil || msg.Reply == "" {
		return nats.ErrInvalidMsg
	}

	var ro broker.ReplyOptions
	for _, opt := range opts {
		opt(&ro)
	}

	reply := &nats.Msg{Subject: msg.Reply, Data: data}
	if len(ro.Header) > 0 {
		reply.Header = toNatsHeader(ro.Header)
	}
	return b.nc.PublishMsg(reply)
}

func (b *Broker) Drain() error {
	if b.nc == nil {
		return nil
	}
	return b.nc.Drain()
}

func (b *Broker) Close() {
	if b.nc == nil {
		return
	}
	b.nc.Close()
}

type subscription struct {
	sub *nats.Subscription
}

func (s *subscription) Unsubscribe() error { return s.sub.Unsubscribe() }
func (s *subscription) Drain() error       { return s.sub.Drain() }

func buildNatsOptions(o *options) []nats.Option {
	opts := []nats.Option{
		nats.Name(o.name),
		nats.Timeout(o.connectTimeout),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(o.reconnectWait),
		nats.MaxReconnects(o.maxReconnects),
		nats.PingInterval(o.pingInterval),
		nats.MaxPingsOutstanding(o.maxPingsOutstanding),
	}

	if o.token != "" {
		opts = append(opts, nats.Token(o.token))
	} else if o.user != "" || o.password != "" {
		opts = append(opts, nats.UserInfo(o.user, o.password))
	}
	if o.tlsCfg != nil {
		opts = append(opts, nats.Secure(o.tlsCfg))
	}
	if len(o.natsOptions) > 0 {
		opts = append(opts, o.natsOptions...)
	}
	return opts
}

func toNatsHeader(h broker.Header) nats.Header {
	if len(h) == 0 {
		return nil
	}
	nh := nats.Header{}
	for k, vals := range h {
		if len(vals) == 0 {
			continue
		}
		// copy to avoid aliasing slices passed by caller
		cp := make([]string, len(vals))
		copy(cp, vals)
		nh[k] = cp
	}
	return nh
}

func fromNatsHeader(h nats.Header) broker.Header {
	if len(h) == 0 {
		return nil
	}
	bh := broker.Header{}
	for k, vals := range h {
		if len(vals) == 0 {
			continue
		}
		cp := make([]string, len(vals))
		copy(cp, vals)
		bh[k] = cp
	}
	return bh
}
