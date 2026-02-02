package gate

import (
	"context"
	"errors"
	"net/url"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/server"
)

var (
	ErrAppNotFound       = errors.New("app not found")
	ErrNetServerRequired = errors.New("net server required")
	ErrLocatorRequired   = errors.New("locator required")
	ErrBrokerRequired    = errors.New("broker required")
	ErrDiscoveryRequired = errors.New("discovery required")
)

type Gate struct {
	opts  *options
	appID string
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
	app, ok := wukong.FromContext(ctx)
	if !ok {
		return ErrAppNotFound
	}
	o := g.opts
	if o.netServer == nil {
		return ErrNetServerRequired
	}
	if o.locator == nil {
		return ErrLocatorRequired
	}
	if o.broker == nil {
		return ErrBrokerRequired
	}
	if o.discovery == nil {
		return ErrDiscoveryRequired
	}
	g.appID = app.ID()
	o.netServer.Start()
	return nil
}

func (g *Gate) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Gate) Endpoint() (*url.URL, error) {
	//TODO implement me
	panic("implement me")
}
