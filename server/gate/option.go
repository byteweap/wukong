package gate

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/network"
)

// options 选项
type options struct {
	netServer network.Server  // 网络服务器
	locator   locator.Locator // 玩家位置定位器
	broker    broker.Broker   // 消息传输代理
}

type Option func(*options)

func defaultOptions() *options {
	return &options{}
}

// NetServer 设置网络服务器
func (g *Gate) NetServer(netServer network.Server) Option {
	return func(o *options) {
		if netServer != nil {
			o.netServer = netServer
		}
	}
}

// Locator 设置玩家位置定位器
func (g *Gate) Locator(locator locator.Locator) Option {
	return func(o *options) {
		if locator != nil {
			o.locator = locator
		}
	}
}

// Broker 设置消息传输代理
func (g *Gate) Broker(broker broker.Broker) Option {
	return func(o *options) {
		if broker != nil {
			o.broker = broker
		}
	}
}
