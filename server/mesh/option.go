package mesh

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
)

const (
	defaultPrefix            = "wukong"
	defaultMessageBufferSize = 256
)

// options 选项
type options struct {
	prefix            string          // subject \ redis key 前缀
	messageBufferSize int             // 消息缓冲区大小
	locator           locator.Locator // 玩家位置定位器
	broker            broker.Broker   // 消息传输代理
}

// Option 定义 Mesh 可选配置函数
type Option func(*options)

// defaultOptions 返回默认配置
func defaultOptions() *options {
	return &options{
		prefix:            defaultPrefix,
		messageBufferSize: defaultMessageBufferSize,
	}
}

// Prefix 设置 broker subject 前缀, 默认: wukong
func Prefix(prefix string) Option {
	return func(o *options) {
		if prefix != "" {
			o.prefix = prefix
		}
	}
}

// MessageBufferSize 设置消息缓冲区大小, 默认 256
func MessageBufferSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.messageBufferSize = size
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
