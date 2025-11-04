package gate

import (
	"fmt"

	"github.com/byteweap/wukong/pkg/knet"
	"github.com/byteweap/wukong/pkg/knet/websocket"
)

type Gate struct {
	opts      *Options
	netServer knet.Server
}

func New(opts ...Option) *Gate {

	// options
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Gate{
		opts: options,
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
	g.netServer.Shutdown()
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
		fmt.Println("Gate start success, addr: ", addr, ", pattern: ", pattern)
	})

	ws.OnStop(func(err error) {
		fmt.Println("Gate websocket stopped, error: ", err)
	})
	ws.OnConnect(func(conn knet.Conn) {
		fmt.Println("Gate connect success, id: ", conn.ID(), ", localAddr: ", conn.LocalAddr(), ", remoteAddr: ", conn.RemoteAddr())
	})
	ws.OnDisconnect(func(conn knet.Conn) {
		fmt.Println("Gate disconnect success, id: ", conn.ID(), ", localAddr: ", conn.LocalAddr(), ", remoteAddr: ", conn.RemoteAddr())
	})
	ws.OnTextMessage(func(conn knet.Conn, msg []byte) {
		fmt.Println("Gate receive text message: ", string(msg))
	})
	ws.OnBinaryMessage(func(conn knet.Conn, msg []byte) {
		fmt.Println("Gate receive binary message: ", msg)
	})
	ws.OnError(func(err error) {
		fmt.Println("Gate err: ", err)
	})
	g.netServer = ws
}
