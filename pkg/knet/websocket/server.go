package websocket

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

// Server websocket server
// 两种方式二选一,使用示例参考 websocket_test
//  1. HandleRequest()
//  2. Run()
type Server struct {
	opts       *Options
	hub        *hub
	httpServer *http.Server
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)
	return &Server{
		opts: options,
		hub:  newHub(options),
	}
}

// HandleRequest 处理请求
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	netConn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("[WebSocket] UpgradeHTTP error: %v", err)
		return
	}
	// allocate connection
	if err = s.hub.allocate(netConn); err != nil {
		log.Printf("[WebSocket] allocate error: %v", err)
		return
	}
}

// Run 运行
func (s *Server) Run() {

	ln, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		log.Fatalf("[WebSocket] net.Listen error: %v", err)
		return
	}
	s.httpServer = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == s.opts.Pattern {
				s.HandleRequest(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	log.Printf("[WebSocket] http server listening at %v%s\n", ln.Addr(), s.opts.Pattern)

	if s.opts.CertFile != "" && s.opts.KeyFile != "" {
		err = s.httpServer.ServeTLS(ln, s.opts.CertFile, s.opts.KeyFile)
	} else {
		err = s.httpServer.Serve(ln)
	}
	if err != nil {
		log.Printf("[WebSocket] http Serve error: %v\n", err)
	}
}

// Shutdown 优雅关闭
func (s *Server) Shutdown() error {

	// http
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	// hub
	if err := s.hub.shutdown(); err != nil {
		return err
	}
	return nil
}
