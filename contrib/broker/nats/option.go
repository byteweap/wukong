package nats

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	// default values (偏稳定与低延迟的折中)
	defaultURLs                = nats.DefaultURL
	defaultName                = "wk-nats-broker"
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

// URLs 设置 NATS 服务地址
func URLs(urls ...string) Option {
	return func(o *options) {
		if len(urls) > 0 {
			o.urls = strings.Join(urls, ",")
		}
	}
}

// Name 设置连接名称
func Name(name string) Option {
	return func(o *options) {
		if name != "" {
			o.name = name
		}
	}
}

// Token 使用 token 认证
func Token(token string) Option {
	return func(o *options) {
		o.token = token
	}
}

// UserPass 使用用户名/密码认证
func UserPass(user, pass string) Option {
	return func(o *options) {
		o.user = user
		o.password = pass
	}
}

// ConnectTimeout 设置连接超时. 默认 3 秒
func ConnectTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.connectTimeout = d
		}
	}
}

// Reconnect 设置重连策略. 默认 250 毫秒, 无限重连
func Reconnect(wait time.Duration, max int) Option {
	return func(o *options) {
		if wait > 0 {
			o.reconnectWait = wait
		}
		o.maxReconnects = max
	}
}

// Ping 设置心跳参数. 默认 20 秒, 3 个心跳未响应则认为连接异常
func Ping(interval time.Duration, maxOutstanding int) Option {
	return func(o *options) {
		if interval > 0 {
			o.pingInterval = interval
		}
		if maxOutstanding > 0 {
			o.maxPingsOutstanding = maxOutstanding
		}
	}
}
