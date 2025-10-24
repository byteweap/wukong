package websocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
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
	state      atomic.Value
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)
	return &Server{
		opts: options,
		hub:  newHub(options),
	}
}

// OnStart 监听服务启动
func (s *Server) OnStart(handler StartHandler) {
	s.opts.startHandler = handler
}

// OnStop 监听服务停止
func (s *Server) OnStop(handler StopHandler) {
	s.opts.stopHandler = handler
}

// OnConnect 监听建立链接
func (s *Server) OnConnect(handler ConnectHandler) {
	s.opts.connectHandler = handler
}

// OnDisconnect 监听断开链接
func (s *Server) OnDisconnect(handler DisconnectHandler) {
	s.opts.disconnectHandler = handler
}

// OnMessage 监听文本消息
func (s *Server) OnMessage(handler MessageHandler) {
	s.opts.messageHandler = handler
}

// OnBinaryMessage 监听二进制消息
func (s *Server) OnBinaryMessage(handler MessageHandler) {
	s.opts.binaryMessageHandler = handler
}

// ErrorHandler 错误处理函数
func (s *Server) ErrorHandler(handler ErrorHandler) {
	s.opts.errorHandler = handler
}

// HandleRequest 处理请求
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	netConn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		s.opts.errorHandler(fmt.Errorf("upgrade http error: %v", err))
		return
	}
	// allocate connection
	if err = s.hub.allocate(netConn); err != nil {
		s.opts.handleError(fmt.Errorf("allocate error: %v", err))
		return
	}
}

// Run 运行
func (s *Server) Run() {

	ln, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		s.opts.errorHandler(fmt.Errorf("net listen error: %v", err))
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

	s.opts.handleStart(s.opts.Addr, s.opts.Pattern)

	if s.opts.CertFile != "" && s.opts.KeyFile != "" {
		err = s.httpServer.ServeTLS(ln, s.opts.CertFile, s.opts.KeyFile)
	} else {
		err = s.httpServer.Serve(ln)
	}
	if err != nil {
		s.opts.handleError(fmt.Errorf("http serve error: %v", err))
	}

	s.opts.handleStop()

}

// Shutdown 优雅关闭
func (s *Server) Shutdown() error {

	// http
	if s.httpServer == nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
	}

	// hub
	if s.hub != nil {
		if err := s.hub.shutdown(); err != nil {
			return err
		}
	}
	return nil
}
