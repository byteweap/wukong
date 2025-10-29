package gate

import (
	"fmt"

	"github.com/byteweap/wukong/pkg/wnet/websocket"
)

type Gate struct {
	opts *Options
	ws   *websocket.Server
}

func New(opts ...Option) *Gate {

	// options
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// websocket server
	ws := websocket.NewServer(
		websocket.WithMaxMessageSize(options.MaxMessageSize),
		websocket.WithMaxConnections(options.MaxConnections),
		websocket.WithReadTimeout(options.ReadTimeout),
		websocket.WithWriteTimeout(options.WriteTimeout),
		websocket.WithWriteQueueSize(options.WriteQueueSize),
	)
	ws.OnStart(func(addr, pattern string) {
		fmt.Println("Gate start success, addr: ", addr, ", pattern: ", pattern)
	})
	ws.OnStop(func() {
		fmt.Println("Gate stop success")
	})
	ws.OnConnect(func(conn *websocket.Conn) {
		fmt.Println("Gate connect success, id: ", conn.ID(), ", localAddr: ", conn.LocalAddr(), ", remoteAddr: ", conn.RemoteAddr())
	})
	ws.OnDisconnect(func(conn *websocket.Conn) {
		fmt.Println("Gate disconnect success, id: ", conn.ID(), ", localAddr: ", conn.LocalAddr(), ", remoteAddr: ", conn.RemoteAddr())
	})
	ws.OnMessage(func(conn *websocket.Conn, msg []byte) {
		fmt.Println("Gate receive message: ", string(msg))
	})
	ws.ErrorHandler(func(err error) {
		fmt.Println("Gate err: ", err)
	})
	return &Gate{
		opts: options,
		ws:   ws,
	}
}

// Start gate server
func (g *Gate) Start() {
	g.ws.Run()
}

// Stop gate server
func (g *Gate) Stop() {
	g.ws.Shutdown()
}
