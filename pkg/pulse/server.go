package pulse

import (
	"context"
	"net/http"

	"github.com/gobwas/ws"
)

type Server struct {
	opts *serverOptions
	hub  *serverHub
}

func NewServer(opts ...ServerOption) *Server {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}
	return &Server{
		opts: o,
		hub:  newServerHub(),
	}
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) error {

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
func (s *Server) BroadcastBinary(msg []byte, filters ...func(conn *ServerConn) bool) {
	s.hub.broadcastBinary(msg, filters...)
}

// BroadcastText 广播文本消息
// filters: 条件过滤器, 返回true则发送消息,否则不发
func (s *Server) BroadcastText(msg []byte, filters ...func(conn *ServerConn) bool) {
	s.hub.broadcastText(msg, filters...)
}

func (s *Server) Close() error {
	s.hub.close()
	return nil
}
