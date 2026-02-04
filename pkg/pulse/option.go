package pulse

import "time"

type (
	OnConnectHandler    func(*Conn)
	OnDisconnectHandler func(*Conn, error)
	OnMessageHandler    func(*Conn, []byte)
	OnErrorHandler      func(*Conn, error)
)

type BackpressureMode int

const (
	BackpressureKick  BackpressureMode = iota // 队列满直接断开（网关常用）
	BackpressureDrop                          // 队列满丢消息（适合低价值同步）
	BackpressureBlock                         // 队列满阻塞写入（不推荐网关）
)

type options struct {
	SendQueueSize  int
	MaxMessageSize int
	ReadTimeout    time.Duration // 0 表示不设置
	WriteTimeout   time.Duration // 0 表示不设置
	Backpressure   BackpressureMode

	// Upgrade 校验：可选
	CheckOrigin func(origin string) bool

	onConnect    OnConnectHandler
	onDisconnect OnDisconnectHandler
	onMessage    OnMessageHandler
	onError      OnErrorHandler
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		SendQueueSize:  256,
		MaxMessageSize: 64 * 1024,
		ReadTimeout:    0,
		WriteTimeout:   0,
		Backpressure:   BackpressureKick,
	}
}

func SendQueueSize(size int) Option {
	return func(o *options) {
		o.SendQueueSize = size
	}
}

func MaxMessageSize(size int) Option {
	return func(o *options) {
		o.MaxMessageSize = size
	}
}
func ReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.WriteTimeout = timeout
	}
}

func Backpressure(mode BackpressureMode) Option {
	return func(o *options) {
		o.Backpressure = mode
	}
}

func CheckOrigin(check func(origin string) bool) Option {
	return func(o *options) {
		o.CheckOrigin = check
	}
}

func OnConnect(fn OnConnectHandler) Option {
	return func(o *options) {
		o.onConnect = fn
	}
}

func OnDisconnect(fn OnDisconnectHandler) Option {
	return func(o *options) {
		o.onDisconnect = fn
	}
}

func OnMessage(fn OnMessageHandler) Option {
	return func(o *options) {
		o.onMessage = fn
	}
}

func OnError(fn OnErrorHandler) Option {
	return func(o *options) {
		o.onError = fn
	}
}
