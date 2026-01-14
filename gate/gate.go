package gate

import (
	"fmt"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/logger"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/contrib/broker/nats"
	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/contrib/logger/zerolog"
	"github.com/byteweap/wukong/contrib/network/websocket"
)

// Gate WebSocket 网关服务器
type Gate struct {
	opts           *Options        // 配置选项
	logger         logger.Logger   // 日志记录器
	netServer      network.Server  // 网络服务器（WebSocket/KCP/TCP）
	sessionManager *SessionManager // 会话管理器
	locator        locator.Locator // 玩家位置定位器
	broker         broker.Broker   // 消息传输代理
}

// New 创建新的网关服务器实例
func New(opts ...Option) (*Gate, error) {
	// 应用配置选项
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	logger := zerolog.New()

	// 初始化定位器和消息代理
	redisOpts, locatorOpts, brokerOpts := options.RedisOptions, options.LocatorOptions, options.BrokerOptions

	locator := redis.New(redisOpts, locatorOpts.KeyFormat, locatorOpts.GateFieldName, locatorOpts.GameFieldName)

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
		opts:           options,
		sessionManager: NewSessionManager(),
		locator:        locator,
		broker:         broker,
	}, nil
}

// Start 启动网关服务器
func (g *Gate) Start() {
	// 初始化网络配置
	g.setupNetwork()
	// 启动网络服务器
	g.netServer.Start()
}

// Stop 停止网关服务器
func (g *Gate) Stop() {
	g.netServer.Stop()
}

// setupNetwork 初始化网络服务器配置
func (g *Gate) setupNetwork() {
	options := g.opts.NetworkOptions

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
