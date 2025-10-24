package gate

import (
	"net/http"
	"strconv"

	"github.com/byteweap/wukong/pkg/knet/websocket"
)

type Gate struct {
	opts Options
	ws   *websocket.Server
}

func New(opts ...Option) *Gate {
	// options
	options := newOptions(opts...)

	// websocket server
	ws := websocket.NewServer(
		websocket.WithMaxMessageSize(options.MaxMessageSize),
		websocket.WithMaxConnections(options.MaxConnections),
		websocket.WithReadTimeout(options.ReadTimeout),
		websocket.WithWriteTimeout(options.WriteTimeout),
		websocket.WithWriteQueueSize(options.WriteQueueSize),
	)

	return &Gate{
		opts: options,
		ws:   ws,
	}
}

func (g *Gate) Run() {
	http.HandleFunc("/ws", g.ws.HandleRequest)
	if err := http.ListenAndServe(g.opts.Addr+":"+strconv.Itoa(g.opts.Port), nil); err != nil {
		panic(err)
	}

}
