package websocket

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/pkg/kos"

	"github.com/gobwas/ws"
)

// Server 定义WebSocket服务器的增强接口
type Server interface {
	network.Server
	// HandleRequest 处理HTTP升级到WebSocket的请求
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

// server 实现Server接口的内部结构体
// 管理WebSocket服务器的生命周期和连接池
type server struct {
	opts       *Options      // 服务器配置选项
	hub        *hub          // 连接管理器
	httpServer *http.Server  // HTTP服务器实例
	state      atomic.Value  // 服务器状态
	done       chan struct{} // 用于通知Start方法执行关闭流程
	closed     int32         // 标记是否已关闭，防止重复关闭
}

// 确保server实现了Server接口
var _ Server = (*server)(nil)

// NewServer 创建新的WebSocket服务器实例
func NewServer(opts ...Option) Server {
	options := newOptions(opts...)
	return &server{
		opts: options,
		hub:  newHub(options),
		done: make(chan struct{}),
	}
}

// Addr 返回服务器监听地址
func (s *server) Addr() string {
	return s.opts.Addr
}

// Protocol 返回协议名称
func (s *server) Protocol() string {
	return "websocket"
}

// OnStart 设置服务启动时的回调函数
func (s *server) OnStart(handler network.StartHandler) {
	s.opts.startHandler = handler
}

// OnStop 设置服务停止时的回调函数
func (s *server) OnStop(handler network.StopHandler) {
	s.opts.stopHandler = handler
}

// OnConnect 设置新连接建立时的回调函数
func (s *server) OnConnect(handler network.ConnectHandler) {
	s.opts.connectHandler = handler
}

// OnDisconnect 设置连接断开时的回调函数
func (s *server) OnDisconnect(handler network.ConnectHandler) {
	s.opts.disconnectHandler = handler
}

// OnTextMessage 设置接收到文本消息时的回调函数
func (s *server) OnTextMessage(handler network.ConnMessageHandler) {
	s.opts.messageHandler = handler
}

// OnBinaryMessage 设置接收到二进制消息时的回调函数
func (s *server) OnBinaryMessage(handler network.ConnMessageHandler) {
	s.opts.binaryMessageHandler = handler
}

// OnError 设置错误发生时的回调函数
func (s *server) OnError(handler network.ErrorHandler) {
	s.opts.errorHandler = handler
}

// HandleRequest 处理HTTP升级到WebSocket的请求
func (s *server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// 将HTTP连接升级为WebSocket连接
	netConn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		s.opts.handleError(fmt.Errorf("upgrade http error: %v", err))
		return
	}
	// 分配连接到hub进行管理
	if err = s.hub.allocate(netConn); err != nil {
		s.opts.handleError(fmt.Errorf("allocate error: %v", err))
		return
	}
}

// Start 启动WebSocket服务器
func (s *server) Start() {
	// 创建TCP监听器
	ln, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		s.opts.handleError(fmt.Errorf("net listen error: %v", err))
		return
	}

	// 配置HTTP服务器
	s.httpServer = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只处理指定路径的请求，其他返回404
			if r.URL.Path == s.opts.Pattern {
				s.HandleRequest(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	// 创建一个通道来接收服务器错误
	serverErr := make(chan error, 1)

	// 启动HTTP服务器
	go func() {
		// 触发启动回调
		s.opts.handleStart(s.opts.Addr, s.opts.Pattern)

		// 根据配置决定是否使用TLS
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

	// 等待退出信号
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
	// 执行关闭流程
	s.shutdown()
}

// Stop 优雅关闭服务器
func (s *server) Stop() {
	// 防止重复关闭
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return
	}
	// 通知Start方法执行关闭流程
	close(s.done)
}

// shutdown 执行服务器关闭的具体流程
// 按顺序关闭HTTP服务器和hub连接管理器
func (s *server) shutdown() {
	// 1. 关闭HTTP服务器，允许30秒的优雅关闭时间
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
	// 3. 触发停止回调
	s.opts.handleStop()
}
