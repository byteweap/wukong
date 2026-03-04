package mesh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/encoding/proto"
	es "github.com/byteweap/wukong/errors"
	"github.com/byteweap/wukong/internal/cluster"
	"github.com/byteweap/wukong/internal/envelope"
	"github.com/byteweap/wukong/pkg/async"
	"github.com/byteweap/wukong/server"
)

type Mesh struct {
	ctx     context.Context
	appID   string // application ID
	appName string // application name

	opts          *options
	routes        sync.Map // key: cmd<<32|version (uint64), value: MessageHandler
	requestRoutes sync.Map // key: cmd.version (string), value: RequestMessageHandler

	onlineHandler    func(uid int64) // 玩家上线
	offlineHandler   func(uid int64) // 玩家掉线
	reconnectHandler func(uid int64) // 玩家重连
}

var _ server.Server = (*Mesh)(nil)

// New 创建 Mesh 服务实例
func New(opts ...Option) *Mesh {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Mesh{opts: o}
}

// Kind 返回服务类型
func (*Mesh) Kind() server.Kind {
	return server.KindMesh
}

// Start 启动 Mesh 服务
func (m *Mesh) Start(ctx context.Context) error {
	app, ok := wukong.FromContext(ctx)
	if !ok {
		return es.ErrAppNotFound
	}

	m.ctx = ctx
	m.appID = app.ID()
	m.appName = app.Name()

	if m.opts.broker == nil {
		return es.ErrBrokerRequired
	}
	if m.opts.locator == nil {
		return es.ErrLocatorRequired
	}
	// 启动常驻协程
	return m.loop()
}

// Stop 停止 Mesh 服务
func (m *Mesh) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// Endpoint 返回服务监听地址
func (m *Mesh) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}

// routeKey 将 cmd/version 打包为 uint64 路由键
func routeKey(cmd, version uint32) uint64 {
	return uint64(cmd)<<32 | uint64(version)
}

func requestRouteKey(cmd, version string) string {
	return cmd + "." + version
}

// OnlineHandler 玩家上线事件处理器
func (m *Mesh) OnlineHandler(handler func(uid int64)) {
	if handler == nil {
		return
	}
	m.onlineHandler = handler
}

// OfflineHandler 玩家掉线事件处理器
func (m *Mesh) OfflineHandler(handler func(uid int64)) {
	if handler == nil {
		return
	}
	m.offlineHandler = handler
}

// ReconnectHandler 玩家重连事件处理器
func (m *Mesh) ReconnectHandler(handler func(uid int64)) {
	if handler == nil {
		return
	}
	m.reconnectHandler = handler
}

// Route 注册业务路由处理器(Gate pub-sub)
// cmd/version 共同确定唯一路由
// handler 支持两种写法，推荐直接传业务函数
//
// 1) 推荐写法
// func(ctx *Context, req *Request)
// 示例: mesh.Route(cmd, version, EnterGame)
//
// 2) 兼容写法
// MessageHandler
// 示例: mesh.Route(cmd, version, mesh.Wrap(EnterGame))
//
// 如果 handler 签名不合法，函数会 panic
func (m *Mesh) Route(cmd, version uint32, handler any) {
	mh, err := adaptMessageHandler(handler)
	if err != nil {
		panic(err)
	}
	key := routeKey(cmd, version)
	m.routes.Store(key, mh)
}

// RequestRoute 注册 request-reply 路由处理器
// cmd/version 共同确定唯一路由
// handler 支持两种写法，推荐直接传业务函数
//
// 1) 推荐写法
// func(ctx *RequestContext, req *Request)
// 示例: mesh.RequestRoute(cmd, version, HandleRequest)
//
// 2) 兼容写法
// RequestMessageHandler
// 示例: mesh.RequestRoute(cmd, version, mesh.WrapRequest(HandleRequest))
//
// 如果 handler 签名不合法，函数会 panic
func (m *Mesh) RequestRoute(cmd, version string, handler any) {
	mh, err := adaptRequestMessageHandler(handler)
	if err != nil {
		panic(err)
	}
	key := requestRouteKey(cmd, version)
	m.requestRoutes.Store(key, mh)
}

// loop 循环
func (m *Mesh) loop() error {

	var (
		o       = m.opts
		subject = cluster.Subject(o.prefix, "*", m.appName, m.appID)
		msgChan = make(chan *broker.Message, o.messageBufferSize)
	)

	// 订阅
	sub, err := o.broker.Sub(m.ctx, subject, func(msg *broker.Message) {
		select {
		case msgChan <- msg:
		case <-m.ctx.Done():
		}
	})
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			// 异常捕获,防止崩溃
			async.Recover(func(r any) {
				log.Errorf("mesh handler panic error: %v", r)
			})
			if err = sub.Close(); err != nil {
				log.Errorf("mesh close subscription error: %v", err)
			}
		}()

		for {
			select {
			case <-m.ctx.Done():
				return
			case msg := <-msgChan:
				if msg == nil {
					continue
				}
				m.handlerMessage(msg)
			}
		}
	}()

	return nil
}

// Request 发送请求并等待回复
func (m *Mesh) Request(subject, cmd, version string, data []byte) (*broker.Message, error) {
	header := broker.Header{}
	header.Set("cmd", cmd)
	header.Set("version", version)
	return m.opts.broker.Request(m.ctx, subject, data, broker.RequestHeader(header))
}

// okReply 发送成功回复
func (m *Mesh) okReply(reqMsg *broker.Message, data []byte) error {
	if reqMsg == nil {
		return nil
	}
	if reqMsg.Reply == "" {
		return errors.New("reply subject is empty")
	}
	return m.opts.broker.Reply(m.ctx, reqMsg, data)
}

// errReply 发送错误回复
func (m *Mesh) errReply(reqMsg *broker.Message, err string) error {
	if reqMsg == nil {
		return errors.New("request message is nil")
	}
	if reqMsg.Reply == "" {
		return errors.New("reply subject is empty")
	}
	header := reqMsg.Header
	header.Set("error", err)
	return m.opts.broker.Reply(m.ctx, reqMsg, nil, broker.ReplyHeader(header))
}

// handlerRequestReplyMessage 来自其它服务的(request-reply)消息
func (m *Mesh) handlerRequestReplyMessage(msg *broker.Message) {
	if msg == nil {
		return
	}
	header := msg.Header
	if header == nil {
		if err := m.errReply(msg, "header is nil"); err != nil {
			log.Errorf("mesh [handlerRequestReplyMessage] reply error: %v", err)
		}
		return
	}
	cmd, version := header.Get("cmd"), header.Get("version")
	if handler, ok := m.requestRoutes.Load(requestRouteKey(cmd, version)); ok {
		handler.(RequestMessageHandler)(m, msg)
	} else {
		if err := m.errReply(msg, fmt.Sprintf("cmd:%s version:%s not found", cmd, version)); err != nil {
			log.Errorf("mesh [handlerRequestReplyMessage] reply error: %v", err)
		}
	}
}

// handlerPubSubMessage 来自Gate的(pub-sub)消息
func (m *Mesh) handlerPubSubMessage(msg *broker.Message) {
	envy := &envelope.Gate2MeshEnvelope{}
	if err := proto.Unmarshal(msg.Data, envy); err != nil {
		log.Errorf("mesh unmarshal Gate2MeshEnvelope error: %v", err)
		return
	}
	meta := envy.GetMeta()
	if meta == nil {
		log.Errorf("mesh missing meta in Gate2MeshEnvelope")
		return
	}
	if handler, ok := m.routes.Load(routeKey(meta.GetCmd(), meta.GetVersion())); ok {
		handler.(MessageHandler)(m, msg, envy)
	}
}

// 处理消息
func (m *Mesh) handlerMessage(msg *broker.Message) {
	if msg.Reply != "" {
		m.handlerRequestReplyMessage(msg)
	} else {
		m.handlerPubSubMessage(msg)
	}
}
