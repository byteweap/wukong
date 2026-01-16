package gate

import (
	"context"
	"time"

	"github.com/byteweap/wukong/log"
	"github.com/redis/go-redis/v9"
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

	// NetworkOptions 网络选项
	NetworkOptions struct {
		Addr           string
		Pattern        string
		MaxMessageSize int64
		MaxConnections int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		WriteQueueSize int
	}

	// LocatorOptions 定位器选项
	LocatorOptions struct {
		KeyFormat     string
		GateFieldName string
		GameFieldName string
	}

	// BrokerOptions 消息代理选项
	BrokerOptions struct {
		Name                string        // 名称, 默认 "gate"
		URLs                []string      // 连接地址,
		Token               string        // 认证 token
		User                string        // 认证用户名
		Password            string        // 认证密码
		ConnectTimeout      time.Duration // 连接超时时间, 默认 3 秒
		ReconnectWait       time.Duration // 重连等待时间, 默认 250 毫秒, 无限重连
		MaxReconnects       int           // 最大重连次数, 默认 0, 无限重连
		PingInterval        time.Duration // 心跳间隔时间, 默认 20 秒, 3 个心跳未响应则认为连接异常
		MaxPingsOutstanding int           // 最大心跳未响应次数, 默认 0, 无限重连
	}

	// RegistryOptions 注册选项
	RegistryOptions struct {
		RegistryTimeout time.Duration // 注册/注销 超时时间, 默认 10 秒
	}

	// options 选项
	options struct {
		ctx    context.Context
		logger log.Logger

		application ApplicationOptions
		network     NetworkOptions
		locator     LocatorOptions
		redis       redis.UniversalOptions
		broker      BrokerOptions
		registry    RegistryOptions
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
		network: NetworkOptions{
			Addr:           "0.0.0.0:9000",
			Pattern:        "/",
			MaxConnections: 10000,
			MaxMessageSize: 4 * 1024, // 4KB
			ReadTimeout:    0,
			WriteTimeout:   0,
			WriteQueueSize: 0,
		},
		locator: LocatorOptions{
			KeyFormat:     "gate:%d",
			GateFieldName: "gate",
			GameFieldName: "game",
		},
		redis: redis.UniversalOptions{
			Addrs:      []string{"localhost:6379"},
			Username:   "",
			Password:   "",
			DB:         0,
			ClientName: "wukong-gate",
		},
		broker: BrokerOptions{
			Name:                "gate",
			URLs:                []string{"localhost:4222"},
			Token:               "",
			User:                "",
			Password:            "",
			ConnectTimeout:      3 * time.Second,
			ReconnectWait:       250 * time.Millisecond,
			MaxReconnects:       -1,
			PingInterval:        20 * time.Second,
			MaxPingsOutstanding: 3,
		},
		registry: RegistryOptions{
			RegistryTimeout: 10 * time.Second,
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

func Addr(addr string) Option {
	return func(o *options) {
		o.application.Addr = addr
		o.network.Addr = addr
	}
}

func Metadata(metadata map[string]string) Option {
	return func(o *options) {
		o.application.Metadata = metadata
	}
}

func Pattern(pattern string) Option {
	return func(o *options) {
		if o.network.Pattern != "" {
			o.network.Pattern = pattern
		}
	}
}

func MaxConnections(maxConnections int) Option {
	return func(o *options) {
		if o.network.MaxConnections > 0 {
			o.network.MaxConnections = maxConnections
		}
	}
}

func MaxMessageSize(maxMessageSize int64) Option {
	return func(o *options) {
		if o.network.MaxMessageSize > 0 {
			o.network.MaxMessageSize = maxMessageSize
		}
	}
}

func ReadTimeout(readTimeout time.Duration) Option {
	return func(o *options) {
		if readTimeout > 0 {
			o.network.ReadTimeout = readTimeout
		}
	}
}

func WriteTimeout(writeTimeout time.Duration) Option {
	return func(o *options) {
		if writeTimeout > 0 {
			o.network.WriteTimeout = writeTimeout
		}
	}
}

func WriteQueueSize(writeQueueSize int) Option {
	return func(o *options) {
		if writeQueueSize > 0 {
			o.network.WriteQueueSize = writeQueueSize
		}
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func Locator(keyFormat, gateFieldName, gameFieldName string) Option {
	return func(o *options) {
		o.locator.KeyFormat = keyFormat
		o.locator.GateFieldName = gateFieldName
		o.locator.GameFieldName = gameFieldName
	}
}

func Redis(opts redis.UniversalOptions) Option {
	return func(o *options) {
		o.redis = opts
	}
}

func Broker(opts BrokerOptions) Option {
	return func(o *options) {
		o.broker = opts
	}
}
