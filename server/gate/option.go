package gate

import (
	"context"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/network"
)

// options 选项
type options struct {
	ctx       context.Context
	logger    log.Logger
	netServer network.Server  // 网络服务器
	locator   locator.Locator // 玩家位置定位器
	broker    broker.Broker   // 消息传输代理
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		ctx: context.Background(),
	}
}
