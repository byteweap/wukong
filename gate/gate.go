package gate

import (
	"context"
	"fmt"
	"sync"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/logger"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/contrib/broker/nats"
	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/contrib/logger/zerolog"
	"github.com/byteweap/wukong/contrib/network/websocket"
	"github.com/google/uuid"
)

// Gate websocket 网关
type Gate struct {
	ctx    context.Context
	cancel context.CancelFunc

	opts           *options          // 配置选项
	logger         logger.Logger     // 日志
	netServer      network.Server    // 网络服务器（WebSocket）
	locator        locator.Locator   // 玩家位置定位器
	broker         broker.Broker     // 消息传输代理
	sessionManager *SessionManager   // 会话管理器
	registry       registry.Registry // 服务注册与发现器

	mu       sync.Mutex
	instance *registry.ServiceInstance // 服务实例
}

// New 创建新的网关服务器实例
func New(opts ...Option) (*Gate, error) {

	// 应用配置选项
	o := defaultOptions()

	if id, err := uuid.NewUUID(); err == nil {
		o.application.ID = id.String()
	}
	for _, opt := range opts {
		opt(o)
	}

	// 选项
	var (
		redisOpts   = o.redis
		locatorOpts = o.locator
		brokerOpts  = o.broker
	)

	// logger
	logger := zerolog.New()

	// 定位器
	locator := redis.New(
		redisOpts,
		locatorOpts.KeyFormat,
		locatorOpts.GateFieldName,
		locatorOpts.GameFieldName,
	)

	// 消息代理
	broker, err := nats.New(
		nats.Name(brokerOpts.Name),
		nats.URLs(brokerOpts.URLs...),
		nats.Token(brokerOpts.Token),
		nats.UserPass(brokerOpts.User, brokerOpts.Password),
		nats.ConnectTimeout(brokerOpts.ConnectTimeout),
		nats.Reconnect(brokerOpts.ReconnectWait, brokerOpts.MaxReconnects),
		nats.Ping(brokerOpts.PingInterval, brokerOpts.MaxPingsOutstanding),
	)
	if err != nil {
		return nil, err
	}

	return &Gate{
		logger:         logger.With("module", "gate"),
		opts:           o,
		sessionManager: NewSessionManager(),
		locator:        locator,
		broker:         broker,
	}, nil
}

// Run 启动网关服务器
func (g *Gate) Run() error {

	instance, err := g.buildInstance()
	if err != nil {
		return err
	}
	g.mu.Lock()
	g.instance = instance
	g.mu.Unlock()

	// 初始化网络配置
	g.setupNetwork()

	// 启动网络服务器
	g.netServer.Start()

	return nil
}

// Stop 关闭网关服务器
func (g *Gate) Stop() {

	g.netServer.Stop()

	if err := g.broker.Close(); err != nil {
		g.logger.Error().Err(err).Msg("broker close error")
	}

	if err := g.locator.Close(); err != nil {
		g.logger.Error().Err(err).Msg("locator close error")
	}

	if err := g.sessionManager.Close(); err != nil {
		g.logger.Error().Err(err).Msg("session manager close error")
	}

	g.logger.Info().Msg("gate server shutdown success")
}

// setupNetwork 初始化网络服务器配置
func (g *Gate) setupNetwork() {
	options := g.opts.network

	// 创建 WebSocket 服务器
	ws := websocket.NewServer(
		websocket.Addr(options.Addr),
		websocket.Pattern(options.Pattern),
		websocket.MaxMessageSize(options.MaxMessageSize),
		websocket.MaxConnections(options.MaxConnections),
		websocket.ReadTimeout(options.ReadTimeout),
		websocket.WriteTimeout(options.WriteTimeout),
		websocket.WriteQueueSize(options.WriteQueueSize),
	)
	ws.OnStart(func(addr, pattern string) {
		g.logger.Info().Msgf("websocket server start success, listening on %s%s", addr, pattern)
	})

	ws.OnStop(func() {
		g.logger.Info().Msg("websocket server stop success")
	})
	ws.OnConnect(func(conn network.Conn) {
		g.logger.Info().Msgf("connect success, id: %d, localAddr: %s, remoteAddr: %s", conn.ID(), conn.LocalAddr(), conn.RemoteAddr())
	})
	ws.OnDisconnect(func(conn network.Conn) {
		g.logger.Info().Msgf("disconnect success, id: %d, localAddr: %s, remoteAddr: %s", conn.ID(), conn.LocalAddr(), conn.RemoteAddr())
	})
	ws.OnBinaryMessage(func(conn network.Conn, msg []byte) {
		g.handlerBinaryMessage(conn, msg)
	})
	ws.OnError(func(err error) {
		g.logger.Error().Err(err).Msg("websocket server err")
	})
	g.netServer = ws
}

// handlerBinaryMessage 处理接收到的二进制消息
func (g *Gate) handlerBinaryMessage(_ network.Conn, msg []byte) {
	fmt.Println("Gate receive binary message: ", msg)
}

// buildInstance 构建服务实例
func (g *Gate) buildInstance() (*registry.ServiceInstance, error) {

	return &registry.ServiceInstance{
		ID:        g.opts.application.ID,
		Name:      g.opts.application.Name,
		Version:   g.opts.application.Version,
		Metadata:  g.opts.application.Metadata,
		Endpoints: []string{g.opts.application.Addr},
	}, nil
}

// registerService 注册服务
func (g *Gate) registerService() error {
	if g.registry != nil && g.instance != nil {
		return g.registry.Register(g.ctx, g.instance)
	}
	return nil
}

// unregisterService 注销服务
func (g *Gate) unregisterService() error {
	if g.registry != nil && g.instance != nil {
		return g.registry.Deregister(g.ctx, g.instance)
	}
	return nil
}
