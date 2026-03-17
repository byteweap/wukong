package gate

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/olahol/melody"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/encoding/proto"
	es "github.com/byteweap/wukong/errors"
	"github.com/byteweap/wukong/internal/cluster"
	"github.com/byteweap/wukong/internal/envelope"
	"github.com/byteweap/wukong/pkg/async"
	"github.com/byteweap/wukong/pkg/endpoint"
	"github.com/byteweap/wukong/pkg/host"
	"github.com/byteweap/wukong/server"
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

func (g *Gate) Kind() server.Kind {
	return server.KindGate
}

func (g *Gate) validate() error {
	if g.opts == nil {
		return es.ErrOptionsRequired
	}
	if g.opts.locator == nil {
		return es.ErrLocatorRequired
	}
	if g.opts.broker == nil {
		return es.ErrBrokerRequired
	}
	if g.opts.discovery == nil {
		return es.ErrDiscoveryRequired
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
		return es.ErrAppNotFound
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
	var e2 error
	if g.ws != nil {
		e2 = g.ws.Close()
	}

	if err := errors.Join(e1, e2); err != nil {
		return err
	}
	log.Info("[gate] server stopped")
	return nil
}

// Endpoint 获取网关地址
func (g *Gate) Endpoint(_ context.Context) (*url.URL, error) {
	if err := g.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return g.endpoint, nil
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

	loc := g.opts.locator

	// 绑定网关
	if err := loc.Bind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] new connection success, bind gate error, uid: %v, err: %v", uid, err)
		return
	}

	// 广播 上线、重连 事件到上游服务
	event := cluster.Event_Online
	if ok {
		event = cluster.Event_Reconnect
	}
	g.broadcastEvent(uid, event)
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
	if err := g.opts.locator.UnBind(g.ctx, uid, g.appName, g.appID); err != nil {
		log.Errorf("[websocket] connection disconnect success, unbind gate error, uid: %v, err: %v", uid, err)
	}

	// 广播掉线事件到上游服务
	g.broadcastEvent(uid, cluster.Event_Offline)
}

// 接收到文本消息时调用
func (g *Gate) onTextMessage(s *melody.Session, msg []byte) {
	// todo
}

// 接收到二进制消息时调用
func (g *Gate) onBinaryMessage(s *melody.Session, msg []byte) {

	meta := &envelope.IMessage{}
	if err := proto.Unmarshal(msg, meta); err != nil {
		log.Errorf("[websocket] unmarshal envelope error: %v", err)
		return
	}
	uids, ok := s.Get("uid")
	if !ok {
		log.Error("[websocket] onBinaryMessage get uid error, session not contains uid key")
		return
	}
	uid := uids.(int64)

	// 业务消息分发
	g.dispatch(uid, meta)
}

// 业务消息分发至 mesh
func (g *Gate) dispatch(uid int64, e *envelope.IMessage) {

	if e == nil {
		log.Errorf("[websocket] dispatch error, envelope is nil")
		return
	}
	var (
		toService   = e.GetService()
		loc, bro, _ = g.opts.locator, g.opts.broker, g.opts.discovery
	)

	curNode, err := loc.Node(g.ctx, uid, toService)
	if err != nil {
		log.Errorf("[websocket] dispatch get mesh node error, uid: %v, toService: %v, err: %v", uid, toService, err)
		return
	}

	data, err := proto.Marshal(e)
	if err != nil {
		log.Errorf("[websocket] dispatch marshal to mesh data error: %v", err)
		return
	}
	node := curNode
	if curNode == "" {
		// todo 确定一个mesh节点(负载均衡算法)
		//services, err := disc.GetService(g.ctx, toService)
		//if err != nil {
		//	log.Errorf("[websocket] dispatch get all services error, uid: %v, app: %v, err: %v", uid, toService, err)
		//	return
		//}
	}
	// 构建消息头
	var (
		reply  = g.Subject(toService) // 回复主题
		header = cluster.BuildHeader(uid, cluster.Event_Business, reply, g.appName, toService)
	)
	// 发布消息到 Mesh
	subject := cluster.Subject(g.opts.prefix, g.appName, toService, node)
	if err = bro.Pub(g.ctx, subject, data, broker.PubHeader(header)); err != nil {
		log.Errorf("[websocket] dispatch error, uid: %v, subject: %v, err: %v", uid, subject, err)
		return
	}
	log.Debugf("[websocket] dispatch success, uid: %v, subject: %v", uid, subject)
}

// 广播系统事件
func (g *Gate) broadcastEvent(uid int64, event cluster.Event) {

	// 获取玩家当前所在所有节点
	snMap, err := g.opts.locator.AllNodes(g.ctx, uid)
	if err != nil {
		log.Errorf("[websocket] broadcast event, get all nodes error, uid: %v, event: %v, err: %v", uid, event, err)
		return
	}
	for service, node := range snMap {
		if service == g.appName { // 不包括 gate
			continue
		}
		// 发布消息到 Mesh
		var (
			header  = cluster.BuildHeader(uid, event, g.Subject(service), g.appName, service)
			subject = cluster.Subject(g.opts.prefix, g.appName, service, node)
		)
		if err = g.opts.broker.Pub(g.ctx, subject, nil, broker.PubHeader(header)); err != nil {
			log.Errorf("[websocket] broadcast event error, uid: %v, subject: %v, err: %v", uid, subject, err)
			return
		}
		log.Debugf("[websocket] broadcast event success, uid: %v, subject: %v, event: %v", uid, subject, event)
	}
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
		select {
		case msgChan <- msg:
		case <-g.ctx.Done():
		}
	})
	if err != nil {
		return err
	}

	// 处理收到的消息
	go func(ctx context.Context, sub broker.Subscription, ch <-chan *broker.Message) {
		defer func() {
			// 异常捕获,防止崩溃
			async.Recover(func(r any) {
				log.Errorf("gate handler panic error: %v", r)
			})
			if err = sub.Close(); err != nil {
				log.Errorf("gate close subscription error: %v", err)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				if msg == nil {
					continue
				}
				g.handleMessage(msg)
			}
		}
	}(g.ctx, sub, msgChan)

	return nil
}

// handlerRequestReplyMessage 来自其它服务的(request-reply)消息
func (g *Gate) handleRequestReplyMessage(msg *broker.Message) {
	// todo
}

// handlerPubSubMessage 来自Mesh服务的(pub-sub)消息
func (g *Gate) handlePubSubMessage(msg *broker.Message) {
	// 2. 直接回复给玩家的消息
	uid := cluster.GetUidBy(msg.Header)
	if uid <= 0 {
		log.Errorf("[websocket] reply2player get uid error, uid: %v", uid)
		return
	}
	session, ok := g.sessions.get(uid)
	if !ok {
		log.Errorf("[websocket] reply2player get session error, uid: %v", uid)
		return
	}
	if err := session.WriteBinary(msg.Data); err != nil {
		log.Errorf("[websocket] reply2player write binary error, uid: %v, err: %v", uid, err)
		return
	}
	log.Debugf("[websocket] reply2player success, uid: %v", uid)
}

// 处理来自其它服务的消息
func (g *Gate) handleMessage(msg *broker.Message) {
	if msg.Reply != "" {
		g.handleRequestReplyMessage(msg)
	} else {
		g.handlePubSubMessage(msg)
	}
}

// Subject 获取当前服务主题
func (g *Gate) Subject(fromApp string) string {
	return cluster.Subject(g.opts.prefix, fromApp, g.appName, g.appID)
}
