package gate

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/olahol/melody"

	"github.com/byteweap/meta"
	"github.com/byteweap/meta/component/broker"
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/component/registry"
	"github.com/byteweap/meta/component/selector"
	es "github.com/byteweap/meta/errors"
	"github.com/byteweap/meta/internal/cluster"
	"github.com/byteweap/meta/pkg/async"
	"github.com/byteweap/meta/pkg/endpoint"
	"github.com/byteweap/meta/pkg/host"
	"github.com/byteweap/meta/server"
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

	mu        sync.RWMutex
	selectors map[string]selector.Selector // 服务节点选择器 key: 服务名
	watchers  map[string]registry.Watcher  // 服务节点监听器 key: 服务名
}

var _ server.Server = (*Gate)(nil)

func New(opts ...Option) *Gate {

	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Gate{
		ctx:       context.Background(),
		opts:      o,
		Server:    &http.Server{},
		sessions:  newSessions(),
		selectors: make(map[string]selector.Selector),
		watchers:  make(map[string]registry.Watcher),
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
	if g.opts.selectorFunc == nil {
		return es.ErrSelectorRequired
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

	m.HandleConnect(g.handleConnect)
	m.HandleDisconnect(g.handleDisconnect)
	m.HandleMessage(g.handleTextMessage)
	m.HandleMessageBinary(g.handleBinaryMessage)
	m.HandleError(g.handleError)
	m.HandleClose(g.handleClose)

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

	app, ok := meta.FromContext(ctx)
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

	err := errors.Join(e1, e2)

	// 3. 停止监听器
	g.mu.Lock()
	for _, watcher := range g.watchers {
		if e := watcher.Stop(); e != nil {
			err = errors.Join(err, e)
		}
	}
	g.mu.Unlock()

	if err != nil {
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

// Subject 获取当前服务主题
func (g *Gate) Subject(fromApp string) string {
	return cluster.Subject(g.opts.prefix, fromApp, g.appName, g.appID)
}

// 确保选择器
func (g *Gate) ensure(service string) (selector.Selector, error) {

	g.mu.RLock()
	sel := g.selectors[service]
	g.mu.RUnlock()
	if sel != nil {
		return sel, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	// double check
	if sel = g.selectors[service]; sel != nil {
		return sel, nil
	}

	sel = g.opts.selectorFunc()

	w, err := g.opts.discovery.Watch(g.ctx, service)
	if err != nil {
		log.Errorf("[websocket] ensure watch service error, service: %v, err: %v", service, err)
		return nil, err
	}
	g.selectors[service] = sel
	g.watchers[service] = w

	go func() {
		for {
			select {
			case <-g.ctx.Done():
				return
			default:
				instances, err := w.Next()
				if err != nil {
					return
				}
				nodes := make([]selector.Node, 0, len(instances))
				for _, instance := range instances {
					nodes = append(nodes,
						selector.NewNode(
							instance.ID,
							instance.Name,
							instance.Version,
							instance.Metadata,
						),
					)
				}
				sel.Update(nodes) // 更新节点
			}
		}
	}()
	return sel, nil
}
