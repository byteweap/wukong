package gate

import (
	"context"
	"time"

	"github.com/byteweap/wukong/component/broker"
	"github.com/byteweap/wukong/component/locator"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/network"
	"github.com/byteweap/wukong/component/registry"
	"github.com/google/uuid"
)

type (
	// Application 应用选项
	Application struct {
		id       string
		name     string
		version  string
		metadata map[string]string
		addr     string
	}

	// options 选项
	options struct {
		ctx             context.Context
		app             Application
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
		app: Application{
			id:       uuid.New().String(),
			name:     "wukong-gate",
			metadata: make(map[string]string),
			addr:     "0.0.0.0:9000",
		},
	}
}

// Context 设置上下文, 默认值: context.Background()
func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// ID 设置服务ID, 默认值: uuid()
func ID(id string) Option {
	return func(o *options) {
		o.app.id = id
	}
}

// Name 设置服务名称, 默认: wukong-gate
func Name(name string) Option {
	return func(o *options) {
		o.app.name = name
	}
}

// Version 设置服务版本, 默认: v1.0.0
func Version(version string) Option {
	return func(o *options) {
		o.app.version = version
	}
}

// Metadata 设置自定义服务元数据
func Metadata(metadata map[string]string) Option {
	return func(o *options) {
		o.app.metadata = metadata
	}
}

// Logger 设置日志记录器, 默认: std logger
func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// NetServer 设置网络服务器, required
func NetServer(netServer network.Server) Option {
	return func(o *options) {
		o.netServer = netServer
	}
}

// Locator 设置玩家位置定位器, required
func Locator(locator locator.Locator) Option {
	return func(o *options) {
		o.locator = locator
	}
}

// Broker 设置消息传输代理, required
func Broker(broker broker.Broker) Option {
	return func(o *options) {
		o.broker = broker
	}
}

// Registry 设置服务注册与发现器、注册超时时间, required
func Registry(registry registry.Registry, registryTimeout time.Duration) Option {
	return func(o *options) {
		o.registry = registry
		o.registryTimeout = registryTimeout
	}
}
