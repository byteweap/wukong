package nats

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/byteweap/wukong/component/broker"
)

// ID 是 NATS broker 实现标识符。
const ID = "nats(core)"

// Broker 使用 NATS Core 实现 broker.Broker。
type Broker struct {
	opts *options
	nc   *nats.Conn
}

var _ broker.Broker = (*Broker)(nil)

// New 创建 NATS Core broker 并立即连接.
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

// NewWith 使用现有的 *nats.Conn 创建 broker(调用者负责其生命周期).
func NewWith(nc *nats.Conn) *Broker {
	return &Broker{opts: defaultOptions(), nc: nc}
}

// ID 返回实现标识.
func (b *Broker) ID() string { return ID }

// Pub 发布一条消息(fire-and-forget).
func (b *Broker) Pub(ctx context.Context, subject string, data []byte, opts ...broker.PublishOption) error {
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

// Sub 订阅主题. 支持队列组订阅，上下文取消时自动退订.
func (b *Broker) Sub(ctx context.Context, subject string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscription, error) {
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

// Request 发送请求并等待响应(request-reply).
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

// Reply 回复一条请求消息.
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

// Close 立即关闭连接.
func (b *Broker) Close() error {
	if b.nc == nil {
		return nil
	}
	if err := b.nc.Drain(); err != nil {
		b.nc.Close()
		return err
	}
	return nil
}

type subscription struct {
	sub *nats.Subscription
}

// Unsub 立即取消订阅.
func (s *subscription) Unsub() error { return s.sub.Unsubscribe() }

// Shutdown 优雅关闭订阅: 完成处理中的消息后退出.
func (s *subscription) Shutdown() error { return s.sub.Drain() }

// buildNatsOptions 根据配置构建 NATS 连接选项。
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

// toNatsHeader 将 broker.Header 转换为 nats.Header.
func toNatsHeader(h broker.Header) nats.Header {
	if len(h) == 0 {
		return nil
	}
	nh := nats.Header{}
	for k, vals := range h {
		if len(vals) == 0 {
			continue
		}
		// 复制以避免与调用者传入的切片共享底层数组
		cp := make([]string, len(vals))
		copy(cp, vals)
		nh[k] = cp
	}
	return nh
}

// fromNatsHeader 将 nats.Header 转换为 broker.Header.
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
