package etcd

import (
	"time"
)

const (
	// 默认值
	defaultNamespace         = "/registry"
	defaultDialTimeout       = 3 * time.Second
	defaultTTL               = 30 * time.Second
	defaultKeepAliveInterval = 10 * time.Second
)

// EtcdConfig etcd 配置
type EtcdConfig struct {
	// Addrs etcd 集群地址列表，格式: "ip:port"
	Addrs []string
	// DialTimeout 连接超时时间，默认 3 秒
	DialTimeout time.Duration
	// Username 用户名（可选）
	Username string
	// Password 密码（可选）
	Password string
}

// DefaultEtcdConfig 返回默认的 etcd 配置
func DefaultEtcdConfig() *EtcdConfig {
	return &EtcdConfig{
		Addrs:       []string{"localhost:2379"},
		DialTimeout: defaultDialTimeout,
	}
}

func validate(cfg *EtcdConfig) {
	if len(cfg.Addrs) == 0 {
		cfg.Addrs = []string{"localhost:2379"}
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = defaultDialTimeout
	}
}

// options 客户端配置选项
type options struct {
	namespace         string
	ttl               time.Duration
	keepAliveInterval time.Duration
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		namespace:         defaultNamespace,
		ttl:               defaultTTL,
		keepAliveInterval: defaultKeepAliveInterval,
	}
}

// Namespace 设置命名空间前缀，默认 "/registry"
func Namespace(ns string) Option {
	return func(o *options) {
		if ns != "" {
			o.namespace = ns
		}
	}
}

// TTL 设置租约 TTL，默认 30 秒
func TTL(ttl time.Duration) Option {
	return func(o *options) {
		if ttl > 0 {
			o.ttl = ttl
		}
	}
}

// KeepAliveInterval 设置续租间隔，默认 10 秒
func KeepAliveInterval(interval time.Duration) Option {
	return func(o *options) {
		if interval > 0 {
			o.keepAliveInterval = interval
		}
	}
}
