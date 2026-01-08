package nats

import (
	"crypto/tls"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	// default values (偏稳定与低延迟的折中)
	defaultURLs                = nats.DefaultURL
	defaultName                = "wukong-broker"
	defaultConnectTimeout      = 3 * time.Second
	defaultReconnectWait       = 250 * time.Millisecond
	defaultMaxReconnects       = -1 // 无限重连
	defaultPingInterval        = 20 * time.Second
	defaultMaxPingsOutstanding = 3
)

type options struct {
	urls string
	name string

	// auth
	token    string
	user     string
	password string
	tlsCfg   *tls.Config

	// connect/reconnect
	connectTimeout      time.Duration
	reconnectWait       time.Duration
	maxReconnects       int
	pingInterval        time.Duration
	maxPingsOutstanding int

	// advanced
	natsOptions []nats.Option
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		urls:                defaultURLs,
		name:                defaultName,
		connectTimeout:      defaultConnectTimeout,
		reconnectWait:       defaultReconnectWait,
		maxReconnects:       defaultMaxReconnects,
		pingInterval:        defaultPingInterval,
		maxPingsOutstanding: defaultMaxPingsOutstanding,
	}
}

// WithURLs 设置 NATS 服务地址（可传逗号分隔的多个 URL）。
func WithURLs(urls string) Option {
	return func(o *options) {
		if urls != "" {
			o.urls = urls
		}
	}
}

// WithName 设置连接名称。
func WithName(name string) Option {
	return func(o *options) {
		if name != "" {
			o.name = name
		}
	}
}

// WithToken 使用 token 认证。
func WithToken(token string) Option {
	return func(o *options) {
		o.token = token
	}
}

// WithUserPass 使用用户名/密码认证。
func WithUserPass(user, pass string) Option {
	return func(o *options) {
		o.user = user
		o.password = pass
	}
}

// WithTLSConfig 设置 TLS 配置。
func WithTLSConfig(cfg *tls.Config) Option {
	return func(o *options) {
		o.tlsCfg = cfg
	}
}

// WithConnectTimeout 设置连接超时。
func WithConnectTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.connectTimeout = d
		}
	}
}

// WithReconnect 设置重连策略。
func WithReconnect(wait time.Duration, max int) Option {
	return func(o *options) {
		if wait > 0 {
			o.reconnectWait = wait
		}
		o.maxReconnects = max
	}
}

// WithPing 设置心跳参数。
func WithPing(interval time.Duration, maxOutstanding int) Option {
	return func(o *options) {
		if interval > 0 {
			o.pingInterval = interval
		}
		if maxOutstanding > 0 {
			o.maxPingsOutstanding = maxOutstanding
		}
	}
}

// WithNatsOptions 透传底层 nats.Option（高级用法）。
func WithNatsOptions(opts ...nats.Option) Option {
	return func(o *options) {
		if len(opts) > 0 {
			o.natsOptions = append(o.natsOptions, opts...)
		}
	}
}
