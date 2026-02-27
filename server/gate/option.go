package gate

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/registry"
)

// options 选项
type options struct {

	// websocket
	path string // ws 路径
	addr string // ws 地址

	locator   locator.Locator   // 玩家位置定位器
	broker    broker.Broker     // 消息传输代理
	discovery registry.Registry // 服务发现
}

type Option func(*options)

func defaultOptions() *options {
	return &options{}
}

// Websocket 设置 websocket 选项
func Websocket(addr, path string) Option {
	return func(o *options) {
		if addr != "" {
			o.addr = addr
		}
		if path != "" {
			o.path = path
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
