package gate

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/pkg/endpoint"
	"github.com/byteweap/wukong/pkg/host"
	"github.com/byteweap/wukong/server"
	"github.com/olahol/melody"
)

var (
	ErrAppNotFound       = errors.New("app not found")
	ErrLocatorRequired   = errors.New("locator required")
	ErrBrokerRequired    = errors.New("broker required")
	ErrDiscoveryRequired = errors.New("discovery required")
)

type Gate struct {
	*http.Server

	opts     *options
	appID    string
	ws       *melody.Melody
	ln       net.Listener
	endpoint *url.URL
}

var _ server.Server = (*Gate)(nil)

func New(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Gate{opts: o, Server: &http.Server{}}
}

func (g *Gate) Kind() server.Kind {
	return server.KindGate
}

func (g *Gate) validate() error {
	if g.opts == nil {
		return errors.New("options required")
	}
	//if g.opts.locator == nil {
	//	return ErrLocatorRequired
	//}
	//if g.opts.broker == nil {
	//	return ErrBrokerRequired
	//}
	//if g.opts.discovery == nil {
	//	return ErrDiscoveryRequired
	//}
	return nil
}

func (g *Gate) setup(appID string) {

	g.appID = appID

	m := melody.New()
	m.HandleConnect(func(s *melody.Session) {
		// todo
	})
	m.HandleDisconnect(func(s *melody.Session) {
		// todo
	})
	//m.HandleMessage(func(s *melody.Session, msg []byte) {
	//	// todo
	//})
	m.HandleMessageBinary(func(s *melody.Session, msg []byte) {
		// todo
	})
	m.HandleError(func(s *melody.Session, err error) {
		// todo
	})
	m.HandleClose(func(s *melody.Session, code int, reason string) error {
		// todo
		return nil
	})

	g.ws = m

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

	if err := g.listenAndEndpoint(); err != nil {
		return err
	}

	g.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != g.opts.path {
			http.NotFound(w, r)
			return
		}
		_ = g.ws.HandleRequest(w, r)
	})

	log.Infof("[websocket] server listening on: %s", g.ln.Addr().String())
	return g.Serve(g.ln)
}

func (g *Gate) Stop(ctx context.Context) error {

	log.Info("[websocket] server stopping")
	err := g.Shutdown(ctx)
	if err != nil {
		if ctx.Err() != nil {
			log.Warn("[websocket] server couldn't stop gracefully in time, doing force stop")
			err = g.Close()
		}
	}
	return err
}

func (g *Gate) listenAndEndpoint() error {
	if g.ln == nil {
		ln, err := net.Listen("tcp", g.opts.addr)
		if err != nil {
			return err
		}
		g.ln = ln
	}
	if g.endpoint == nil {
		addr, err := host.Extract(g.opts.addr, g.ln)
		if err != nil {
			return err
		}
		g.endpoint = endpoint.NewEndpoint(endpoint.Scheme("ws", false), addr)
	}
	return nil
}

func (g *Gate) Endpoint() (*url.URL, error) {
	if err := g.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return g.endpoint, nil
}

// 消息分发至 Game
func (g *Gate) dispatch() {

}
