package pulse

import (
	"time"
)

type (
	ConnectHandler    func(*Conn)
	DisconnectHandler func(*Conn, error)
	MessageHandler    func(*Conn, []byte)
	ErrorHandler      func(*Conn, error)
)

type BackpressureMode int

const (
	BackpressureKick  BackpressureMode = iota // 队列满直接断开（网关常用）
	BackpressureDrop                          // 队列满丢消息（适合低价值同步）
	BackpressureBlock                         // 队列满阻塞写入（不推荐网关）
)

type options struct {
	// 发送队列大小
	sendQueueSize int
	// 最大消息大小
	maxMessageSize int64
	// 读超时时间, 0 表示不设置
	readTimeout time.Duration
	// 写超时时间, 0 表示不设置
	writeTimeout time.Duration
	// 背压模式
	backpressure BackpressureMode

	// Upgrade 校验：可选
	checkOrigin func(origin string) bool

	// 连接建立回调
	onConnect ConnectHandler
	// 连接断开回调
	onDisconnect DisconnectHandler
	// 文本消息回调
	onTextMessage MessageHandler
	// 二进制消息回调
	onBinaryMessage MessageHandler
	// 错误回调
	onError ErrorHandler
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		sendQueueSize:  256,
		maxMessageSize: 64 * 1024,
		readTimeout:    time.Second * 5,
		writeTimeout:   time.Second * 5,
		backpressure:   BackpressureKick,
	}
}

// SendQueueSize 设置发送队列大小, 默认:256
func SendQueueSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.sendQueueSize = size
		}
	}
}

// MaxMessageSize 设置最大消息大小, 默认:64K
func MaxMessageSize(size int64) Option {
	return func(o *options) {
		o.maxMessageSize = size
	}
}

// ReadTimeout 设置读超时, 默认:5秒
func ReadTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.readTimeout = timeout
	}
}

// WriteTimeout 设置写超时, 默认:5秒
func WriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.writeTimeout = timeout
	}
}

// Backpressure 设置背压模式, 默认:BackpressureKick (即队列满直接断开)
func Backpressure(mode BackpressureMode) Option {
	return func(o *options) {
		o.backpressure = mode
	}
}

// CheckOrigin 设置跨域校验函数
func CheckOrigin(check func(origin string) bool) Option {
	return func(o *options) {
		o.checkOrigin = check
	}
}

// OnConnect 设置连接建立时的回调函数
func OnConnect(fn ConnectHandler) Option {
	return func(o *options) {
		o.onConnect = fn
	}
}

func OnDisconnect(fn DisconnectHandler) Option {
	return func(o *options) {
		o.onDisconnect = fn
	}
}

func OnTextMessage(fn MessageHandler) Option {
	return func(o *options) {
		o.onTextMessage = fn
	}
}

func OnBinaryMessage(fn MessageHandler) Option {
	return func(o *options) {
		o.onBinaryMessage = fn
	}
}

func OnError(fn ErrorHandler) Option {
	return func(o *options) {
		o.onError = fn
	}
}
