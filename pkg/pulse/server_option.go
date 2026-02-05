package pulse

import (
	"time"

	"github.com/gobwas/ws"
)

type (
	OnServerConnectHandler    func(*ServerConn)
	OnServerDisconnectHandler func(*ServerConn, error)
	OnServerMessageHandler    func(*ServerConn, ws.OpCode, []byte)
	OnServerErrorHandler      func(*ServerConn, error)
)

type BackpressureMode int

const (
	BackpressureKick  BackpressureMode = iota // 队列满直接断开（网关常用）
	BackpressureDrop                          // 队列满丢消息（适合低价值同步）
	BackpressureBlock                         // 队列满阻塞写入（不推荐网关）
)

type serverOptions struct {
	sendQueueSize  int
	maxMessageSize int64
	readTimeout    time.Duration // 0 表示不设置
	writeTimeout   time.Duration // 0 表示不设置
	backpressure   BackpressureMode

	// Upgrade 校验：可选
	checkOrigin func(origin string) bool

	onConnect    OnServerConnectHandler
	onDisconnect OnServerDisconnectHandler
	onMessage    OnServerMessageHandler
	onError      OnServerErrorHandler
}

type ServerOption func(*serverOptions)

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		sendQueueSize:  256,
		maxMessageSize: 64 * 1024,
		readTimeout:    0,
		writeTimeout:   0,
		backpressure:   BackpressureKick,
	}
}

func ServerSendQueueSize(size int) ServerOption {
	return func(o *serverOptions) {
		o.sendQueueSize = size
	}
}

func ServerMaxMessageSize(size int64) ServerOption {
	return func(o *serverOptions) {
		o.maxMessageSize = size
	}
}
func ServerReadTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.readTimeout = timeout
	}
}

func ServerWriteTimeout(timeout time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.writeTimeout = timeout
	}
}

func ServerBackpressure(mode BackpressureMode) ServerOption {
	return func(o *serverOptions) {
		o.backpressure = mode
	}
}

func ServerCheckOrigin(check func(origin string) bool) ServerOption {
	return func(o *serverOptions) {
		o.checkOrigin = check
	}
}

func OnServerConnect(fn OnServerConnectHandler) ServerOption {
	return func(o *serverOptions) {
		o.onConnect = fn
	}
}

func OnServerDisconnect(fn OnServerDisconnectHandler) ServerOption {
	return func(o *serverOptions) {
		o.onDisconnect = fn
	}
}

func OnServerMessage(fn OnServerMessageHandler) ServerOption {
	return func(o *serverOptions) {
		o.onMessage = fn
	}
}

func OnServerError(fn OnServerErrorHandler) ServerOption {
	return func(o *serverOptions) {
		o.onError = fn
	}
}
