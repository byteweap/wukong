package websocket

import (
	"context"
	"errors"
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
	done       chan struct{} // 用于通知Start方法执行关闭流程
	closed     int32         // 标记是否已关闭，防止重复关闭
}

var _ Server = (*server)(nil)

func NewServer(opts ...Option) Server {
	options := newOptions(opts...)
	return &server{
		opts: options,
		hub:  newHub(options),
		done: make(chan struct{}),
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

	// 创建一个通道来接收服务器错误
	serverErr := make(chan error, 1)

	go func() {
		s.opts.handleStart(s.opts.Addr, s.opts.Pattern)

		if s.opts.CertFile != "" && s.opts.KeyFile != "" {
			err = s.httpServer.ServeTLS(ln, s.opts.CertFile, s.opts.KeyFile)
		} else {
			err = s.httpServer.Serve(ln)
		}

		// 当Serve返回时，通常是因为服务器被关闭
		// 只有在非正常关闭时才报告错误
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("http server error: %v", err)
		}
		close(serverErr)
	}()

	select {
	case <-kos.Signal():
		// 接收到系统信号
	case <-s.done:
		// Stop方法被调用
	case err := <-serverErr:
		// 服务器异常关闭
		if err != nil {
			s.opts.handleError(err)
		}
	}
	s.shutdown()
}

// Stop 优雅关闭
func (s *server) Stop() {
	// 防止重复关闭
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return
	}
	// 通知Start方法执行关闭流程
	close(s.done)
}

func (s *server) shutdown() {

	// 1. 关闭HTTP服务器
	if s.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			s.opts.handleError(fmt.Errorf("http server shutdown error: %v", err))
		}
	}

	// 2. 关闭hub和所有连接
	if s.hub != nil {
		if err := s.hub.shutdown(); err != nil {
			s.opts.handleError(fmt.Errorf("hub shutdown error: %v", err))
		}
	}

	// 优雅关闭服务器
	s.opts.handleStop()
}
