package game

import (
	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
)

// options 选项
type options struct {
	locator locator.Locator // 玩家位置定位器
	broker  broker.Broker   // 消息传输代理
}

type Option func(*options)

func defaultOptions() *options {
	return &options{}
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
