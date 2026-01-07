package gate

import (
	"fmt"

	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/contrib/logger/zerolog"
	"github.com/byteweap/wukong/contrib/network/websocket"
	"github.com/byteweap/wukong/plugin/locator"
	"github.com/byteweap/wukong/plugin/logger"
	"github.com/byteweap/wukong/plugin/network"
)

type Gate struct {
	opts           *Options
	logger         logger.Logger
	netServer      network.Server
	sessionManager *SessionManager
	locator        locator.Locator
}

func New(opts ...Option) (*Gate, error) {

	// options
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	logger := zerolog.New()

	locator := redis.New(options.RedisOptions, options.LocatorKeyFormat, options.LocatorGateFieldName, options.LocatorGameFieldName)

	return &Gate{
		logger:         logger.With("module", "gate"),
		opts:           options,
		sessionManager: NewSessionManager(),
		locator:        locator,
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

	options := g.opts

	// websocket server
	ws := websocket.NewServer(
		websocket.WithAddr(options.Addr),
		websocket.WithPattern(options.Pattern),
		websocket.WithMaxMessageSize(options.MaxMessageSize),
		websocket.WithMaxConnections(options.MaxConnections),
		websocket.WithReadTimeout(options.ReadTimeout),
		websocket.WithWriteTimeout(options.WriteTimeout),
		websocket.WithWriteQueueSize(options.WriteQueueSize),
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
