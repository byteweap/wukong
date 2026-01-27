package gate

import (
	"context"
	"net/url"

	"github.com/byteweap/wukong/server"
)

type Gate struct {
	opts *options
}

var _ server.Server = (*Gate)(nil)

func New(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Gate{opts: o}
}

func (g *Gate) Kind() server.Kind {
	return server.KindGate
}

func (g *Gate) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Gate) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Gate) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}
