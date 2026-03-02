package gate

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/encoding/proto"
	"github.com/byteweap/wukong/internal/cluster"
	"github.com/byteweap/wukong/internal/envelope"
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

	ctx     context.Context
	appID   string // application ID
	appName string // application name

	opts     *options       // server options
	endpoint *url.URL       // server endpoint
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
		ctx:      context.Background(),
		opts:     o,
		Server:   &http.Server{},
		sessions: newSessions(),
	}
}

func (g *Gate) Kind() cluster.Kind {
	return cluster.KindGate
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

func (g *Gate) setup(name, appID string, ctx context.Context) error {

	// base
	g.appID = appID
	g.appName = name
	g.ctx = ctx

	// 监听端口并设置 endpoint
	if err := g.listenAndEndpoint(); err != nil {
		return err
	}

	// websocket
	m, o := melody.New(), g.opts
	m.Config.WriteWait = o.writeTimeout
	m.Config.PongWait = o.pongTimeout
	m.Config.PingPeriod = o.pingInterval
	m.Config.MaxMessageSize = o.maxMessageSize
	m.Config.MessageBufferSize = o.messageBufferSize
	m.Config.ConcurrentMessageHandling = false
	m.HandleConnect(g.onConnect)
	m.HandleDisconnect(g.onDisconnect)
	m.HandleMessage(g.onTextMessage)
	m.HandleMessageBinary(g.onBinaryMessage)
	m.HandleError(func(s *melody.Session, err error) {
		log.Errorf("[websocket] error occurred, err: %v", err)
	})
	m.HandleClose(func(s *melody.Session, code int, reason string) error {
		log.Infof("[websocket] connection closed, code: %v, reason: %v", code, reason)
		return nil
	})
	g.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != g.opts.path {
			http.NotFound(w, r)
			return
		}
		_ = m.HandleRequest(w, r)
	})
	g.ws = m

	return nil
}

// Start 启动网关
func (g *Gate) Start(ctx context.Context) error {

	app, ok := wukong.FromContext(ctx)
	if !ok {
		return ErrAppNotFound
	}
	// 验证参数
	if err := g.validate(); err != nil {
		return err
	}
	// 启动前设置
	if err := g.setup(app.Name(), app.ID(), ctx); err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}
	// 循环(常驻协程)
	if err := g.loop(); err != nil {
		return fmt.Errorf("loop failed: %w", err)
	}

	log.Infof("[gate] server started")
	log.Infof("[websocket] server listening on: %s", g.ln.Addr().String())

	// 启动服务
	return g.Serve(g.ln)
}

// Stop 停止网关
func (g *Gate) Stop(ctx context.Context) error {

	// 1. Shutdown http server
	e1 := g.Shutdown(ctx)
	if e1 != nil && ctx.Err() != nil {
		log.Warn("[websocket] server couldn't stop gracefully in time, doing force stop")
		e1 = g.Close()
	}

	// 2. Close melody
	e2 := g.ws.Close()

	if err := errors.Join(e1, e2); err != nil {
		return err
	}
	log.Info("[gate] server stopped")
	return nil
}

// 监听端口并设置 endpoint
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

// Endpoint 获取网关地址
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
	if err := g.opts.locator.BindGate(g.ctx, uid, g.appID); err != nil {
		log.Errorf("[websocket] new connection success, bind gate error, uid: %v, err: %v", uid, err)
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
	if err := g.opts.locator.UnBindGate(g.ctx, uid, g.appID); err != nil {
		log.Errorf("[websocket] connection disconnect success, unbind gate error, uid: %v, err: %v", uid, err)
	}
}

// 接收到文本消息时调用
func (g *Gate) onTextMessage(s *melody.Session, msg []byte) {
	// todo
}

// 接收到二进制消息时调用
func (g *Gate) onBinaryMessage(s *melody.Session, msg []byte) {

	e := &envelope.Envelope{}
	if err := proto.Unmarshal(msg, e); err != nil {
		log.Errorf("[websocket] unmarshal envelope error: %v", err)
		return
	}
	uids, ok := s.Get("uid")
	if !ok {
		log.Error("[websocket] onBinaryMessage get uid error, session not contains uid key")
		return
	}
	uid := uids.(int64)

	g.dispatch(&envelope.EnvelopeGate2Game{
		E:   e,
		Uid: uid,
	})
}

// 消息分发至 Game
func (g *Gate) dispatch(e *envelope.EnvelopeGate2Game) {

	var (
		uid, toApp  = e.Uid, e.E.App
		loc, bro, _ = g.opts.locator, g.opts.broker, g.opts.discovery
	)

	curGameNode, err := loc.Game(g.ctx, uid)
	if err != nil {
		log.Errorf("[websocket] dispatch get game node error, uid: %v, err: %v", uid, err)
		return
	}
	data, err := proto.Marshal(e)
	if err != nil {
		log.Errorf("[websocket] dispatch marshal to game data error: %v", err)
		return
	}
	node := curGameNode
	if curGameNode == "" {
		// todo 确定一个game节点(负载均衡算法)
		//services, err := disc.GetService(g.ctx, toApp)
		//if err != nil {
		//	log.Errorf("[websocket] dispatch get all services error, uid: %v, app: %v, err: %v", uid, toApp, err)
		//	return
		//}

	}
	// 发布消息到 Game
	subject := cluster.Subject(g.opts.prefix, g.appName, toApp, node)
	if err = bro.Pub(g.ctx, subject, data); err != nil {
		log.Errorf("[websocket] dispatch error, uid: %v, subject: %v, err: %v", uid, subject, err)
		return
	}
	log.Infof("[websocket] dispatch success, uid: %v, subject: %v", uid, subject)
}

// 循环处理来自其它服务的消息
func (g *Gate) loop() error {

	var (
		o       = g.opts
		bro     = g.opts.broker
		msgChan = make(chan *broker.Message, o.messageBufferSize)
		subject = cluster.Subject(o.prefix, "*", g.appName, g.appID)
	)

	// 订阅消息
	sub, err := bro.Sub(g.ctx, subject, func(msg *broker.Message) {
		msgChan <- msg
	})
	if err != nil {
		return err
	}

	// 处理消息
	go func() {
		for {
			select {
			case <-g.ctx.Done():
				_ = sub.Close()
				return
			case msg := <-msgChan:
				log.Info("[websocket] loop receive message", msg)
				if msg.Reply != "" {
					// todo request-reply
				} else {
					// todo pub-sub
				}
			}
		}
	}()

	return nil
}
