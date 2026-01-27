package game

import (
	"context"
	"net/url"

	"github.com/byteweap/wukong/server"
)

type Game struct {
	opts *options
}

var _ server.Server = (*Game)(nil)

func New(opts ...Option) *Game {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Game{opts: o}
}

func (g *Game) Kind() server.Kind {
	return server.KindGame
}

func (g *Game) Start(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Game) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Game) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}
