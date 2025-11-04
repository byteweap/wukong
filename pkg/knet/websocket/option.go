package websocket

import (
	"time"

	"github.com/byteweap/wukong/pkg/knet"
)

// Options 定义WebSocket服务器的配置选项
// 使用选项模式允许灵活配置服务器参数
// 所有参数都有合理的默认值，适合游戏等高频场景
type Options struct {

	// Addr 监听地址
	// 使用 Start() 启动时有效, HandleRequest() 时无效
	Addr string

	// Pattern 监听路径
	// 使用 Start() 启动时有效, HandleRequest() 时无效
	Pattern string

	// 证书文件
	CertFile string

	// 秘钥文件
	KeyFile string

	// MaxMessageSize 最大消息大小
	// 0表示不限制
	// 建议根据应用的消息大小调整此参数，防止内存溢出
	MaxMessageSize int64

	// WriteQueueSize 写入队列的大小
	// 增大队列可以提高写入吞吐量，但会增加内存使用
	// 适合高频、低延迟场景的优化参数
	WriteQueueSize int

	// ReadTimeout 读取超时时间
	// 0表示不超时
	// 建议设置合理的超时，避免空闲连接占用资源
	ReadTimeout time.Duration

	// WriteTimeout 写入超时时间
	// 0表示不超时
	// 建议设置合理的超时，避免网络阻塞
	WriteTimeout time.Duration

	// MaxConnections 最大并发连接数
	// 0表示无限制
	// 限制并发连接数可以防止资源耗尽
	MaxConnections int

	// PingInterval 心跳间隔时间
	// 0表示不启用心跳
	// 建议设置合理的心跳间隔，避免连接断开
	PingInterval time.Duration

	// PongTimeout 心跳响应超时时间
	// 0表示不超时
	// 建议设置合理的超时，避免连接断开
	PongTimeout time.Duration

	// startHandler 服务启动处理器
	startHandler knet.StartHandler

	// stopHandler 服务停止处理器
	stopHandler knet.StopHandler

	// messageHandler 文本消息处理器
	messageHandler knet.ConnMessageHandler

	// binaryMessageHandler 二进制消息处理器
	binaryMessageHandler knet.ConnMessageHandler

	// connectHandler 建立链接处理器
	connectHandler knet.ConnectHandler

	// disconnectHandler 连接断开处理器
	disconnectHandler knet.ConnectHandler

	// errorHandler 错误处理器,可用于打印日志
	errorHandler knet.ErrorHandler
}

// Option 定义选项函数类型
// 用于配置Options结构体
type Option func(*Options)

// newOptions 创建并初始化默认选项
func newOptions(options ...Option) *Options {
	// 设置默认值
	opts := &Options{
		CertFile:       "",               // 证书文件
		KeyFile:        "",               // 密钥文件
		Addr:           ":8000",          // 监听地址
		Pattern:        "/",              // 监听路径
		MaxMessageSize: 1024,             // 默认最大消息大小1KB，适合一般游戏消息
		WriteQueueSize: 1024,             // 默认队列大小为1024，适合中等规模应用
		ReadTimeout:    30 * time.Second, // 默认读超时30秒
		WriteTimeout:   10 * time.Second, // 默认写超时10秒
		MaxConnections: 10000,            // 默认最大连接数10000
		PingInterval:   2 * time.Second,  // 默认心跳间隔30秒
	}

	// 应用用户提供的选项
	for _, option := range options {
		option(opts)
	}

	return opts
}

func (o Options) handleStart(addr, pattern string) {
	if o.startHandler != nil {
		o.startHandler(addr, pattern)
	}
}

func (o Options) handleStop(err error) {
	if o.stopHandler != nil {
		o.stopHandler(err)
	}
}

func (o Options) handleMessage(conn *Conn, message []byte) {
	if o.messageHandler != nil {
		o.messageHandler(conn, message)
	}
}

func (o Options) handleBinaryMessage(conn *Conn, message []byte) {
	if o.binaryMessageHandler != nil {
		o.binaryMessageHandler(conn, message)
	}
}

func (o Options) handleConnect(conn *Conn) {
	if o.connectHandler != nil {
		o.connectHandler(conn)
	}
}

func (o Options) handleDisconnect(conn *Conn) {
	if o.disconnectHandler != nil {
		o.disconnectHandler(conn)
	}
}

func (o Options) handleError(err error) {
	if o.errorHandler != nil {
		o.errorHandler(err)
	}
}

// WithAddr 监听地址
// 默认 :8000, 在使用 Start() 启动时有效
func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

// WithPattern 监听路径
// 默认 /, 在使用 Start() 启动时有效
func WithPattern(pattern string) Option {
	return func(o *Options) {
		o.Pattern = pattern
	}
}

// WithSSL ssl安全配置
// 默认无
func WithSSL(certFile, keyFile string) Option {
	return func(o *Options) {
		o.CertFile = certFile
		o.KeyFile = keyFile
	}
}

// WithMaxMessageSize 最大消息尺寸
// 默认 1KB (1024B)
func WithMaxMessageSize(size int64) Option {
	return func(o *Options) {
		o.MaxMessageSize = size
	}
}

// WithWriteQueueSize 设置写入队列大小
// 默认 1024
func WithWriteQueueSize(size int) Option {
	return func(o *Options) {
		o.WriteQueueSize = size
	}
}

// WithReadTimeout 设置读取超时时间
// 0表示不超时
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = timeout
	}
}

// WithWriteTimeout 设置写入超时时间
// 0表示不超时
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = timeout
	}
}

// WithMaxConnections 设置最大并发连接数
// 0表示无链接数限制
func WithMaxConnections(max int) Option {
	return func(o *Options) {
		o.MaxConnections = max
	}
}

// WithPingInterval 设置心跳间隔时间
// 设置0则默认为2秒
func WithPingInterval(interval time.Duration) Option {
	return func(o *Options) {
		if interval <= 0 {
			interval = 2 * time.Second
		}
		o.PingInterval = interval
	}
}
