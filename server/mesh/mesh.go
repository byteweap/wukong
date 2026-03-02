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
