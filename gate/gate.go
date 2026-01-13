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

// Gate is the websocket gate server.
type Gate struct {
	// opts is the options.
	opts *Options
	// logger is the logger.
	logger logger.Logger
	// netServer is the network server for websocket/kcp/tcp.
	netServer network.Server
	// sessionManager is the session manager.
	sessionManager *SessionManager
	// locator is the locator for player location.
	locator locator.Locator
	// broker is the broker for message transmission.
	broker broker.Broker
}

// New creates a new gate server.
func New(opts ...Option) (*Gate, error) {

	// options
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	logger := zerolog.New()

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

// Start gate server
func (g *Gate) Start() {

	// 1. setup network
	g.setupNetwork()
	// 2. start
	g.netServer.Start()
}

// Stop gate server
func (g *Gate) Stop() {
	g.netServer.Stop()
}

// setupNetwork setup network
func (g *Gate) setupNetwork() {

	options := g.opts.NetworkOptions

	// websocket server
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

// 处理二进制消息
func (g *Gate) handlerBinaryMessage(_ network.Conn, msg []byte) {
	fmt.Println("Gate receive binary message: ", msg)
}
