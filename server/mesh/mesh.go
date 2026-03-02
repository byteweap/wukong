package mesh

import (
	"context"
	"net/url"

	"github.com/byteweap/wukong/server"
)

type Mesh struct {
	opts *options
}

var _ server.Server = (*Mesh)(nil)

func New(opts ...Option) *Mesh {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Mesh{opts: o}
}

func (g *Mesh) Kind() server.Kind {
	return server.KindMesh
}

func (g *Mesh) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Mesh) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Mesh) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}
