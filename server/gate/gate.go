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
	ln net.Listener

	opts     *options // server options
	appID    string   // application ID
	endpoint *url.URL // server endpoint

	ws       *melody.Melody // WebSocket server
	sessions *Sessions      // player sessions
}

var _ server.Server = (*Gate)(nil)

func New(opts ...Option) *Gate {

	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Gate{
		opts:     o,
		Server:   &http.Server{},
		sessions: newSessions(),
	}
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

	m, o := melody.New(), g.opts
	m.Config.WriteWait = o.writeTimeout
	m.Config.PongWait = o.pongTimeout
	//m.Config.PingPeriod = o.pingInterval
	//m.Config.MaxMessageSize = o.maxMessageSize
	//m.Config.MessageBufferSize = o.messageBufferSize
	//m.Config.ConcurrentMessageHandling = false
	m.HandleConnect(g.onConnect)
	m.HandleDisconnect(g.onDisconnect)
	m.HandleMessage(g.onTextMessage)
	m.HandleMessageBinary(g.onBinaryMessage)
	m.HandleError(func(s *melody.Session, err error) {
		// todo
		log.Errorf("[websocket] error occurred, err: %v", err)
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
	// 1. Shutdown http server
	e1 := g.Shutdown(ctx)
	if e1 != nil && ctx.Err() != nil {
		log.Warn("[websocket] server couldn't stop gracefully in time, doing force stop")
		e1 = g.Close()
	}

	// 2. Close melody
	e2 := g.ws.Close()

	return errors.Join(e1, e2)
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

// 连接建立时调用
func (g *Gate) onConnect(s *melody.Session) {

	req, addr := s.Request, s.RemoteAddr()

	uid := g.opts.userIdExtractor(req)
	if uid <= 0 {
		_ = s.Write([]byte("uid is required"))
		_ = s.Close()
		return
	}

	// 注册会话
	session, ok := g.sessions.get(uid)
	if ok {
		log.Warnf("[websocket] connection exists: uid: %v, close old connection", uid)
		_ = session.Close()
	}
	s.Set("uid", uid)
	g.sessions.register(uid, s)

	log.Infof("[websocket] new connection success, uid: %v, %s", uid, addr.String())

	// 绑定网关
	loc := g.opts.locator
	if loc != nil {
		if err := loc.BindGate(context.Background(), uid, g.appID); err != nil {
			log.Errorf("[websocket] new connection success, bind gate error, uid: %v, err: %v", uid, err)
		}
	}
}

// 连接断开时调用
func (g *Gate) onDisconnect(s *melody.Session) {

	uids, ok := s.Get("uid")
	if !ok {
		log.Error("[websocket] connection disconnect error, session not contains uid key")
		return
	}
	uid := uids.(int64)

	// 注销会话
	curSession, yes := g.sessions.get(uid)
	if !yes {
		log.Errorf("[websocket] connection disconnect error, uid: %v not found", uid)
		return
	}
	if curSession != s {
		log.Warnf("[websocket] connection disconnect error, uid: %v session not match", uid)
		return
	}
	g.sessions.unregister(uid)

	log.Infof("[websocket] connection disconnect success, uid: %v", uid)

	// 解绑网关
	loc := g.opts.locator
	if loc != nil {
		if err := loc.UnBindGate(context.Background(), uid, g.appID); err != nil {
			log.Errorf("[websocket] connection disconnect success, unbind gate error, uid: %v, err: %v", uid, err)
		}
	}
}

// 接收到文本消息时调用
func (g *Gate) onTextMessage(s *melody.Session, msg []byte) {
	//todo
}

// 接收到二进制消息时调用
func (g *Gate) onBinaryMessage(s *melody.Session, msg []byte) {
	//todo
}

// 消息分发至 Game
func (g *Gate) dispatch() {
	// todo
}
