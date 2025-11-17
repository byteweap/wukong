package gate

import (
	"fmt"

	"github.com/byteweap/wukong/pkg/klog"
	"github.com/byteweap/wukong/pkg/knet"
	"github.com/byteweap/wukong/pkg/knet/websocket"
)

type Gate struct {
	opts      *Options
	logger    klog.Logger
	netServer knet.Server
}

func New(logger klog.Logger, opts ...Option) *Gate {

	// options
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Gate{
		logger: logger.With("module", "gate"),
		opts:   options,
	}
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
	ws.OnConnect(func(conn knet.Conn) {
		g.logger.Info().Msgf("connect success, id: %d, localAddr: %s, remoteAddr: %s", conn.ID(), conn.LocalAddr(), conn.RemoteAddr())
	})
	ws.OnDisconnect(func(conn knet.Conn) {
		g.logger.Info().Msgf("disconnect success, id: %d, localAddr: %s, remoteAddr: %s", conn.ID(), conn.LocalAddr(), conn.RemoteAddr())
	})
	ws.OnTextMessage(func(conn knet.Conn, msg []byte) {
		fmt.Println("Gate receive text message: ", string(msg))
	})
	ws.OnBinaryMessage(func(conn knet.Conn, msg []byte) {
		fmt.Println("Gate receive binary message: ", msg)
	})
	ws.OnError(func(err error) {
		g.logger.Error().Err(err).Msg("websocket server err")
	})
	g.netServer = ws
}
