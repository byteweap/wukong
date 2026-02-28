package gate

import (
	"net/http"
	"time"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/pkg/conv"
)

// IdExtractor 用户id提取器
// gate 会在建立连接时调用此函数获取用户id
type IdExtractor func(r *http.Request) int64

// options 选项
type options struct {

	// app
	subjectPrefix   string      // 消息主题前缀
	userIdExtractor IdExtractor // 用户 id 提取器

	// websocket
	path              string        // ws 路径
	addr              string        // ws 地址
	writeTimeout      time.Duration // write 超时时间
	pongTimeout       time.Duration // Pong 超时时间
	pingInterval      time.Duration // Ping 间隔时间
	maxMessageSize    int64         // 最大消息大小
	messageBufferSize int           // 消息缓冲区大小, websocket 和 broker 都用

	// component
	locator   locator.Locator   // 玩家位置定位器
	broker    broker.Broker     // 消息传输代理
	discovery registry.Registry // 服务发现
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		subjectPrefix:     "wukong",
		path:              "/",
		addr:              ":9000",
		writeTimeout:      5 * time.Second,
		pongTimeout:       60 * time.Second,
		pingInterval:      10 * time.Second,
		maxMessageSize:    1024 * 2, // 2k
		messageBufferSize: 256,
		userIdExtractor: func(r *http.Request) int64 {
			return conv.Int64(r.FormValue("uid"))
		},
	}
}

// SubjectPrefix 设置消息主题前缀
func SubjectPrefix(prefix string) Option {
	return func(o *options) {
		if o.subjectPrefix != "" {
			o.subjectPrefix = prefix
		}
	}
}

// Addr 设置 websocket 地址
func Addr(addr string) Option {
	return func(o *options) {
		if o.addr != "" {
			o.addr = addr
		}
	}
}

// Path 设置 websocket 路径
func Path(path string) Option {
	return func(o *options) {
		if o.path != "" {
			o.path = path
		}
	}
}

// WriteTimeout 设置 websocket write 超时时间
func WriteTimeout(timeout time.Duration) Option {
	return func(o *options) {
		if o.writeTimeout != 0 {
			o.writeTimeout = timeout
		}
	}
}

// PongTimeout 设置 websocket Pong 超时时间
func PongTimeout(timeout time.Duration) Option {
	return func(o *options) {
		if o.pongTimeout != 0 {
			o.pongTimeout = timeout
		}
	}
}

// PingInterval 设置 websocket Ping 间隔时间
func PingInterval(interval time.Duration) Option {
	return func(o *options) {
		if o.pingInterval != 0 {
			o.pingInterval = interval
		}
	}
}

// MaxMessageSize 设置 websocket 最大消息大小, 单位:字节, 默认: 512
func MaxMessageSize(size int64) Option {
	return func(o *options) {
		if o.maxMessageSize != 0 {
			o.maxMessageSize = size
		}
	}
}

// MessageBufferSize 设置 websocket 消息缓冲区大小,默认: 256
func MessageBufferSize(size int) Option {
	return func(o *options) {
		if o.messageBufferSize != 0 {
			o.messageBufferSize = size
		}
	}
}

// UserIdExtractor 设置用户 id 提取器
// gate 会在建立连接时调用此函数获取用户id, 默认: func(r *http.Request) int64 { return conv.Int64(r.FormValue("uid")) }
func UserIdExtractor(extractor IdExtractor) Option {
	return func(o *options) {
		if extractor != nil {
			o.userIdExtractor = extractor
		}
	}
}

// Locator 设置玩家位置定位器
func Locator(locator locator.Locator) Option {
	return func(o *options) {
		if locator != nil {
			o.locator = locator
		}
	}
}

// Broker 设置消息传输代理
func Broker(broker broker.Broker) Option {
	return func(o *options) {
		if broker != nil {
			o.broker = broker
		}
	}
}

// Discovery 设置服务发现
func Discovery(discovery registry.Registry) Option {
	return func(o *options) {
		if discovery != nil {
			o.discovery = discovery
		}
	}
}
