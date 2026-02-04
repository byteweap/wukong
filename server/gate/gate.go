package gate

import (
	"context"
	"errors"
	"net/url"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/network"
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

func (g *Gate) validate() error {
	if g.opts == nil {
		return errors.New("options required")
	}
	if g.opts.netServer == nil {
		return ErrNetServerRequired
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

	ns := g.opts.netServer

	ns.OnStart(func(addr, pattern string) {
		log.Infof("gate start success, addr: %s, pattern: %s", addr, pattern)
	})
	ns.OnStop(func() {
		log.Infof("gate stop success")
	})
	ns.OnConnect(func(conn network.Conn) {
		log.Infof("gate connection established, remote addr: %s", conn.RemoteAddr())
	})
	ns.OnDisconnect(func(conn network.Conn) {
		log.Infof("gate connection closed, remote addr: %s", conn.RemoteAddr())
	})
	ns.OnError(func(err error) {
		log.Errorf("gate error: %v", err)
	})
	ns.OnTextMessage(func(conn network.Conn, msg []byte) {
		log.Infof("gate received text message, remote addr: %s, message: %s", conn.RemoteAddr(), string(msg))
	})
	ns.OnBinaryMessage(func(conn network.Conn, msg []byte) {
		log.Infof("gate received binary message, remote addr: %s, message: %s", conn.RemoteAddr(), msg)
	})
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

	g.opts.netServer.Start()

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
