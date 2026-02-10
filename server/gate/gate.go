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

func (g *Gate) validate() error {
	if g.opts == nil {
		return errors.New("options required")
	}
	if g.opts.locator == nil {
		return ErrLocatorRequired
	}
	if g.opts.broker == nil {
		return ErrBrokerRequired
	}
	if g.opts.discovery == nil {
		return ErrDiscoveryRequired
	}
	return nil
}

func (g *Gate) setup(appID string) {

	g.appID = appID

}

func (g *Gate) Start(ctx context.Context) error {

	app, ok := wukong.FromContext(ctx)
	if !ok {
		return ErrAppNotFound
	}

	if err := g.validate(); err != nil {
		return err
	}

	g.setup(app.ID())

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

// 消息分发至 Game
func (g *Gate) dispatch() {

}
