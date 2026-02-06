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

func (s *Pulse) HandleRequest(w http.ResponseWriter, r *http.Request) error {

	// Origin 校验
	if s.opts.checkOrigin != nil {
		origin := r.Header.Get("Origin")
		if origin != "" && !s.opts.checkOrigin(origin) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return context.Canceled
		}
	}

	raw, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return err
	}

	return s.hub.allocate(s.opts, raw)
}

// BroadcastBinary 广播二进制消息
// filters: 条件过滤器, 返回true则发送消息,否则不发
func (s *Pulse) BroadcastBinary(msg []byte, filters ...func(conn *Conn) bool) {
	s.hub.broadcastBinary(msg, filters...)
}

// BroadcastText 广播文本消息
// filters: 条件过滤器, 返回true则发送消息,否则不发
func (s *Pulse) BroadcastText(msg []byte, filters ...func(conn *Conn) bool) {
	s.hub.broadcastText(msg, filters...)
}

// NumConnections 返回当前连接数
func (s *Pulse) NumConnections() int {
	return len(s.hub.cs)
}

func (s *Pulse) Close() error {
	s.hub.close()
	return nil
}
