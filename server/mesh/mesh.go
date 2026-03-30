package mesh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/byteweap/meta"
	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/encoding/proto"
	"github.com/byteweap/meta/envelope"
	es "github.com/byteweap/meta/errors"
	"github.com/byteweap/meta/internal/cluster"
	"github.com/byteweap/meta/pkg/async"
	"github.com/byteweap/meta/pkg/conv"
	"github.com/byteweap/meta/server"
)

type Mesh struct {
	ctx     context.Context
	appID   string // application ID
	appName string // application name
	running bool

	opts          *options
	routes        sync.Map // key: cmd<<32|version (uint64), value: MessageHandler
	requestRoutes sync.Map // key: cmd.version (string), value: RpcMessageHandler

	onlineHandler    func(uid int64) // 玩家上线
	offlineHandler   func(uid int64) // 玩家掉线
	reconnectHandler func(uid int64) // 玩家重连

	cancel context.CancelFunc
	done   chan struct{}
	mu     sync.Mutex
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

	if m.opts.broker == nil {
		return es.ErrBrokerRequired
	}
	if m.opts.locator == nil {
		return es.ErrLocatorRequired
	}
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("mesh already running")
	}
	runCtx, cancel := context.WithCancel(ctx)
	m.ctx = runCtx
	m.cancel = cancel
	m.appID = app.ID()
	m.appName = app.Name()
	m.done = make(chan struct{})
	m.running = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		if m.done != nil {
			close(m.done)
			m.done = nil
		}
		m.running = false
		m.cancel = nil
		m.mu.Unlock()
	}()

	// 启动常驻协程
	return m.loop()
}

// Stop 停止 Mesh 服务
func (m *Mesh) Stop(ctx context.Context) error {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return nil
	}
	cancel := m.cancel
	done := m.done
	m.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if done == nil {
		return nil
	}
	select {
	case <-done:
		log.Infof("%s.%s server stop success", m.appName, m.appID)
		return nil
	case <-ctx.Done():
		log.Warnf("%s.%s server stop timeout", m.appName, m.appID)
		return ctx.Err()
	}
}

// Endpoint 返回服务监听地址
func (m *Mesh) Endpoint(ctx context.Context) (*url.URL, error) {
	app, ok := wukong.FromContext(ctx)
	if !ok {
		return nil, es.ErrAppNotFound
	}
	host := app.Name() + "." + app.ID()
	return &url.URL{
		Scheme: "mesh",
		Host:   host + ":0000",
	}, nil
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

// RouteX 注册业务路由处理器(Gate pub-sub)
// cmd/version 共同确定唯一路由
// handler 支持两种写法，推荐直接传业务函数
//
// 1) 推荐写法
// func(ctx *Context, req *Request)
// 示例: mesh.RouteX(cmd, version, EnterGame)
//
// 2) 兼容写法
// MessageHandler
// 示例: mesh.RouteX(cmd, version, mesh.Wrap(EnterGame))
//
// 如果 handler 签名不合法，函数会 panic
// 注意: 使用反射, 热点路由请使用 Route
func (m *Mesh) RouteX(cmd, version uint32, handler any) {
	mh, err := adaptMessageHandler(handler)
	if err != nil {
		panic(err)
	}
	key := routeKey(cmd, version)
	m.routes.Store(key, mh)
}

// Route 注册业务路由处理器
// 该方法要求显式传入 MessageHandler（通常通过 mesh.Wrap 构造）
// 运行期不经过反射调用，适合高频热点路由
func (m *Mesh) Route(cmd, version uint32, handler MessageHandler) {
	if handler == nil {
		panic("mesh: handler is nil")
	}
	key := routeKey(cmd, version)
	m.routes.Store(key, handler)
}

// RpcRouteX 注册 request-reply 路由处理器
// cmd/version 共同确定唯一路由
// handler 支持两种写法，推荐直接传业务函数
//
// 1) 推荐写法
// func(ctx *RpcContext, req *Request) ([]byte, string, int)
// 示例: mesh.RpcRouteX(cmd, version, HandleRequest)
//
// 2) 兼容写法
// RpcMessageHandler
// 示例: mesh.RpcRouteX(cmd, version, mesh.WrapRpc(HandleRequest))
//
// 如果 handler 签名不合法，函数会 panic
// 注意: 使用反射, 热点路由请使用 RpcRoute
func (m *Mesh) RpcRouteX(cmd, version string, handler any) {
	mh, err := adaptRpcMessageHandler(handler)
	if err != nil {
		panic(err)
	}
	key := requestRouteKey(cmd, version)
	m.requestRoutes.Store(key, mh)
}

// RpcRoute 注册 request-reply 路由处理器
// 该方法要求显式传入 RpcMessageHandler（通常通过 mesh.WrapRpc 构造）
// 运行期不经过反射调用，适合高频热点路由
func (m *Mesh) RpcRoute(cmd, version string, handler RpcMessageHandler) {
	if handler == nil {
		panic("mesh: request-reply handler is nil")
	}
	key := requestRouteKey(cmd, version)
	m.requestRoutes.Store(key, handler)
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

	log.Infof("%s.%s server start success", m.appName, m.appID)

	defer func() {
		// 异常捕获,防止崩溃
		async.Recover(func(r any) {
			log.Errorf("mesh handler panic error: %v", r)
		})
		if err := sub.Close(); err != nil {
			log.Errorf("mesh close subscription error: %v", err)
		}
	}()

	for {
		select {
		case <-m.ctx.Done():
			return nil
		case msg := <-msgChan:
			if msg == nil {
				continue
			}
			m.handlerMessage(msg)
		}
	}
}

// Request 发送请求并等待回复
// 入参
//   - subject: 目标服务
//   - cmd: 命令
//   - version: 版本
//   - data: 业务数据
//
// 出参
//   - data: 业务数据
//   - tip: 提示信息
//   - code: 业务状态码
//   - error: 错误信息
func (m *Mesh) Request(subject, cmd, version string, data []byte) ([]byte, string, int, error) {
	header := broker.Header{}
	header.Set("cmd", cmd)
	header.Set("version", version)
	result, err := m.opts.broker.Request(m.ctx, subject, data, broker.RequestHeader(header))
	if err != nil {
		return nil, "", 0, err
	}
	code := result.Header.Get("code")
	tip := result.Header.Get("tip")
	return result.Data, tip, conv.Int(code), nil
}

// sendMessage 发送消息
func (m *Mesh) sendMessage(subject, toService string, bytes []byte, uids ...int64) error {
	var err error
	for _, uid := range uids {
		header := cluster.BuildHeader(uid, cluster.Event_Business, "", m.appName, toService)
		if e := m.opts.broker.Pub(m.ctx, subject, bytes, broker.PubHeader(header)); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}

// okReply 发送成功回复
func (m *Mesh) okReply(reqMsg *broker.Message, data []byte) error {
	if reqMsg == nil {
		return nil
	}
	if reqMsg.Reply == "" {
		return errors.New("reply subject is empty")
	}
	header := reqMsg.Header
	if header == nil {
		header = broker.Header{}
	}
	header.Set("code", "200")
	header.Set("tip", "ok")
	return m.opts.broker.Reply(m.ctx, reqMsg, data, broker.ReplyHeader(header))
}

// errReply 发送错误回复
func (m *Mesh) errReply(reqMsg *broker.Message, code int, tip string) error {
	if reqMsg == nil {
		return errors.New("request message is nil")
	}
	if reqMsg.Reply == "" {
		return errors.New("reply subject is empty")
	}
	header := reqMsg.Header
	if header == nil {
		header = broker.Header{}
	}
	header.Set("code", conv.String(code))
	header.Set("tip", tip)
	return m.opts.broker.Reply(m.ctx, reqMsg, nil, broker.ReplyHeader(header))
}

// handlerRequestReplyMessage 来自其它服务的(request-reply)消息
func (m *Mesh) handlerRequestReplyMessage(msg *broker.Message) {
	if msg == nil {
		return
	}
	header := msg.Header
	if header == nil {
		if err := m.errReply(msg, 100, "header is nil"); err != nil {
			log.Errorf("mesh [handlerRequestReplyMessage] reply error: %v", err)
		}
		return
	}
	cmd, version := header.Get("cmd"), header.Get("version")
	if handler, ok := m.requestRoutes.Load(requestRouteKey(cmd, version)); ok {
		data, tip, code := handler.(RpcMessageHandler)(m, msg)
		m.replyRequestResult(msg, data, tip, code)
	} else {
		if err := m.errReply(msg, 404, fmt.Sprintf("cmd:%s version:%s not found", cmd, version)); err != nil {
			log.Errorf("mesh [handlerRequestReplyMessage] reply error: %v", err)
		}
	}
}

func (m *Mesh) replyRequestResult(msg *broker.Message, data []byte, tip string, code int) {
	if code == http.StatusOK {
		if err := m.okReply(msg, data); err != nil {
			log.Errorf("mesh request-reply ok reply error: %v", err)
		}
		return
	}
	if tip == "" {
		tip = "request failed"
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}
	if err := m.errReply(msg, code, tip); err != nil {
		log.Errorf("mesh request-reply err reply error: %v", err)
	}
}

// handlerPubSubMessage 来自Gate的(pub-sub)消息
func (m *Mesh) handlerPubSubMessage(msg *broker.Message) {

	var (
		uid   = cluster.GetUidBy(msg.Header)
		event = cluster.GetEventBy(msg.Header)
	)

	switch event {
	case cluster.Event_Online:
		if m.onlineHandler != nil {
			m.onlineHandler(uid)
		}
	case cluster.Event_Offline:
		if m.offlineHandler != nil {
			m.offlineHandler(uid)
		}
	case cluster.Event_Reconnect:
		if m.reconnectHandler != nil {
			m.reconnectHandler(uid)
		}
	default:
		e := &envelope.IMessage{}
		if err := proto.Unmarshal(msg.Data, e); err != nil {
			log.Errorf("mesh unmarshal Gate2MeshEnvelope error: %v", err)
			return
		}
		header := e.GetHeader()
		cmd, version := header.GetCmd(), header.GetVersion()
		if handler, ok := m.routes.Load(routeKey(cmd, version)); ok {
			handler.(MessageHandler)(m, msg, e)
		}
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
