package pulse

import (
	"context"
	"net/http"

	"github.com/gobwas/ws"
)

type Pulse struct {
	opts *options
	hub  *hub
}

func New(opts ...Option) *Pulse {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Pulse{
		opts: o,
		hub:  newHub(),
	}
}

func (p *Pulse) HandleRequest(w http.ResponseWriter, r *http.Request) error {

	// Origin 校验
	if p.opts.checkOrigin != nil {
		origin := r.Header.Get("Origin")
		if origin != "" && !p.opts.checkOrigin(origin) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return context.Canceled
		}
	}

	raw, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return err
	}

	p.hub.allocate(context.Background(), p.opts, raw)

	return nil
}

// BroadcastBinary 广播二进制消息
// filters: 条件过滤器, 返回true则发送消息,否则不发
func (p *Pulse) BroadcastBinary(msg []byte, filters ...func(conn *Conn) bool) {
	p.hub.broadcastBinary(msg, filters...)
}

// BroadcastText 广播文本消息
// filters: 条件过滤器, 返回true则发送消息,否则不发
func (p *Pulse) BroadcastText(msg []byte, filters ...func(conn *Conn) bool) {
	p.hub.broadcastText(msg, filters...)
}
