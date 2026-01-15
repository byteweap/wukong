package gate

import (
	"time"

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

	// Options 选项
	Options struct {
		Application ApplicationOptions
		Network     NetworkOptions
		Locator     LocatorOptions
		Redis       redis.UniversalOptions
		Broker      BrokerOptions
	}
)

type Option func(*Options)

func defaultOptions() *Options {

	return &Options{
		Application: ApplicationOptions{
			ID:       "",
			Name:     "gate",
			Version:  "1.0.0",
			Metadata: make(map[string]string),
			Addr:     "0.0.0.0:9000",
		},
		Network: NetworkOptions{
			Addr:           "0.0.0.0:9000",
			Pattern:        "/",
			MaxConnections: 10000,
			MaxMessageSize: 4 * 1024, // 4KB
			ReadTimeout:    0,
			WriteTimeout:   0,
			WriteQueueSize: 0,
		},
		Locator: LocatorOptions{
			KeyFormat:     "gate:%d",
			GateFieldName: "gate",
			GameFieldName: "game",
		},
		Redis: redis.UniversalOptions{
			Addrs:      []string{"localhost:6379"},
			Username:   "",
			Password:   "",
			DB:         0,
			ClientName: "wukong-gate",
		},
		Broker: BrokerOptions{
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
	}
}

func Addr(addr string) Option {
	return func(o *Options) {
		o.Network.Addr = addr
	}
}

func Pattern(pattern string) Option {
	return func(o *Options) {
		if o.Network.Pattern != "" {
			o.Network.Pattern = pattern
		}
	}
}

func MaxConnections(maxConnections int) Option {
	return func(o *Options) {
		if o.Network.MaxConnections > 0 {
			o.Network.MaxConnections = maxConnections
		}
	}
}

func Locator(keyFormat, gateFieldName, gameFieldName string) Option {
	return func(o *Options) {
		o.Locator.KeyFormat = keyFormat
		o.Locator.GateFieldName = gateFieldName
		o.Locator.GameFieldName = gameFieldName
	}
}

func Redis(opts redis.UniversalOptions) Option {
	return func(o *Options) {
		o.Redis = opts
	}
}

func Broker(opts BrokerOptions) Option {
	return func(o *Options) {
		o.Broker = opts
	}
}
