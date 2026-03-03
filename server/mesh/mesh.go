package mesh

import (
	"context"
	"net/url"
	"sync"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/broker"
	es "github.com/byteweap/wukong/errors"
	"github.com/byteweap/wukong/internal/cluster"
	"github.com/byteweap/wukong/server"
)

type Mesh struct {
	ctx     context.Context
	appID   string // application ID
	appName string // application name

	opts   *options
	routes sync.Map // key: route+version, value: MessageHandler

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

// Route 注册业务路由处理器
// cmd/version 共同确定唯一路由
// handler 支持两种写法，推荐直接传业务函数
//
// 1) 推荐写法
// func(ctx *Context, req *Request) error
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

// loop 循环
func (m *Mesh) loop() error {
	var (
		o        = m.opts
		subject  = cluster.Subject(o.prefix, "*", m.appName, m.appID)
		gateChan = make(chan *broker.Message, o.messageBufferSize)
	)
	sub, err := o.broker.Sub(m.ctx, subject, func(msg *broker.Message) {
		select {
		case gateChan <- msg:
		case <-m.ctx.Done():
		}
	})
	if err != nil {
		return err
	}

	go func(ctx context.Context, sub broker.Subscription, ch <-chan *broker.Message) {
		defer sub.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				if msg == nil {
					continue
				}
				m.handlerMessage(msg)
			}
		}
	}(m.ctx, sub, gateChan)

	return nil
}

// 处理 gate 消息
func (m *Mesh) handlerMessage(msg *broker.Message) {
	// todo
}
