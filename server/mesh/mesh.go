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

func (m *Mesh) Route(cmd, version int32, handler MessageHandler) {
	// todo
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
