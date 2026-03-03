package mesh

import (
	"context"
	"net/url"
	"sync"

	"github.com/byteweap/wukong/server"
)

type Mesh struct {
	opts   *options
	routes sync.Map // key: route+version, value: MessageHandler

	onlineHandler    func(uid int64) // 玩家上线
	offlineHandler   func(uid int64) // 玩家掉线
	reconnectHandler func(uid int64) // 玩家重连
}

var _ server.Server = (*Mesh)(nil)

func New(opts ...Option) *Mesh {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Mesh{opts: o}
}

func (*Mesh) Kind() server.Kind {
	return server.KindMesh
}

func (m *Mesh) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *Mesh) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *Mesh) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}

// High 32 bits: cmd, low 32 bits: version.
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

// Route 注册路由
//
// handler 支持两种类型 (Request为自定义类型):
//
//   - MessageHandler mesh.Wrap(func(ctx context.Context, req *Request) error)
//     如: mesh.Route(cmd, version, mesh.Wrap(func(ctx context.Context, req *Request) error))
//
//   - func(ctx context.Context, req *Request) error
//     如: mesh.Route(cmd, version, func(ctx context.Context, req *Request) error)
func (m *Mesh) Route(cmd, version uint32, handler any) {
	mh, err := adaptMessageHandler(handler)
	if err != nil {
		panic(err)
	}
	key := routeKey(cmd, version)
	m.routes.Store(key, mh)
}
