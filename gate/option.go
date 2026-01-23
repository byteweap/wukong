package gate

import (
	"context"
	"time"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/component/registry"
)

type (
	// ApplicationOptions 应用选项
	ApplicationOptions struct {
		ID       string
		Name     string
		Version  string
		Metadata map[string]string
		Addr     string
	}

	// options 选项
	options struct {
		ctx             context.Context
		application     ApplicationOptions
		logger          log.Logger
		netServer       network.Server    // 网络服务器
		locator         locator.Locator   // 玩家位置定位器
		broker          broker.Broker     // 消息传输代理
		registry        registry.Registry // 服务注册与发现器
		registryTimeout time.Duration     // 服务注册与发现超时时间
	}
)

type Option func(*options)

func defaultOptions() *options {

	return &options{
		application: ApplicationOptions{
			ID:       "",
			Name:     "gate",
			Version:  "1.0.0",
			Metadata: make(map[string]string),
			Addr:     "0.0.0.0:9000",
		},
	}
}

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}
func ID(id string) Option {
	return func(o *options) {
		o.application.ID = id
	}
}

func Name(name string) Option {
	return func(o *options) {
		o.application.Name = name
	}
}

func Version(version string) Option {
	return func(o *options) {
		o.application.Version = version
	}
}

func Metadata(metadata map[string]string) Option {
	return func(o *options) {
		o.application.Metadata = metadata
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func NetServer(netServer network.Server) Option {
	return func(o *options) {
		o.netServer = netServer
	}
}

func Locator(locator locator.Locator) Option {
	return func(o *options) {
		o.locator = locator
	}
}

func Broker(broker broker.Broker) Option {
	return func(o *options) {
		o.broker = broker
	}
}

func Registry(registry registry.Registry, registryTimeout time.Duration) Option {
	return func(o *options) {
		o.registry = registry
		o.registryTimeout = registryTimeout
	}
}
