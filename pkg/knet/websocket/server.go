package websocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"

	"github.com/byteweap/wukong/pkg/knet"
	"github.com/byteweap/wukong/pkg/kos"
)

// Server websocket server
type Server interface {
	knet.Server
	HandleRequest(w http.ResponseWriter, r *http.Request)
}
type server struct {
	opts       *Options
	hub        *hub
	httpServer *http.Server
	state      atomic.Value
}

var _ Server = (*server)(nil)

func NewServer(opts ...Option) Server {
	options := newOptions(opts...)
	return &server{
		opts: options,
		hub:  newHub(options),
	}
}

func (s *server) Addr() string {
	return s.opts.Addr
}

func (s *server) Protocol() string {
	return "websocket"
}

// OnStart 监听服务启动
func (s *server) OnStart(handler knet.StartHandler) {
	s.opts.startHandler = handler
}

// OnStop 监听服务停止
func (s *server) OnStop(handler knet.StopHandler) {
	s.opts.stopHandler = handler
}

// OnConnect 监听建立链接
func (s *server) OnConnect(handler knet.ConnectHandler) {
	s.opts.connectHandler = handler
}

// OnDisconnect 监听断开链接
func (s *server) OnDisconnect(handler knet.ConnectHandler) {
	s.opts.disconnectHandler = handler
}

// OnTextMessage 监听文本消息
func (s *server) OnTextMessage(handler knet.ConnMessageHandler) {
	s.opts.messageHandler = handler
}

// OnBinaryMessage 监听二进制消息
func (s *server) OnBinaryMessage(handler knet.ConnMessageHandler) {
	s.opts.binaryMessageHandler = handler
}

// OnError 错误处理函数
func (s *server) OnError(handler knet.ErrorHandler) {
	s.opts.errorHandler = handler
}

// HandleRequest 处理请求
func (s *server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	netConn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		s.opts.handleError(fmt.Errorf("upgrade http error: %v", err))
		return
	}
	// allocate connection
	if err = s.hub.allocate(netConn); err != nil {
		s.opts.handleError(fmt.Errorf("allocate error: %v", err))
		return
	}
}

// Start server
func (s *server) Start() {

	ln, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		s.opts.handleError(fmt.Errorf("net listen error: %v", err))
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

	go func() {
		s.opts.handleStart(s.opts.Addr, s.opts.Pattern)

		if s.opts.CertFile != "" && s.opts.KeyFile != "" {
			err = s.httpServer.ServeTLS(ln, s.opts.CertFile, s.opts.KeyFile)
		} else {
			err = s.httpServer.Serve(ln)
		}
		if err != nil {
			s.opts.handleError(fmt.Errorf("http serve error: %v", err))
		}
	}()

	kos.WaitSignal()

	s.opts.handleStop(nil)
}

// Shutdown 优雅关闭
func (s *server) Shutdown() {

	// http
	if s.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.opts.handleError(fmt.Errorf("http shutdown error: %v", err))
			return
		}
	}

	// hub
	if s.hub != nil {
		if err := s.hub.shutdown(); err != nil {
			s.opts.handleError(fmt.Errorf("hub shutdown error: %v", err))
			return
		}
	}
}
