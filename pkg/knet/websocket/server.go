package websocket

import (
	"log"
	"net/http"

	"github.com/gobwas/ws"
)

type Server struct {
	opts *Options
	hub  *hub
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)
	return &Server{
		opts: options,
		hub:  newHub(options),
	}
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	netConn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("[WebSocket] UpgradeHTTP error: %v", err)
		return
	}
	// 分配一个链接
	if err := s.hub.allocate(netConn); err != nil {
		log.Printf("[WebSocket] allocate error: %v", err)
		return
	}
}
