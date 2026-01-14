package etcd

import (
	"crypto/tls"
	"time"
)

const (
	// 默认值
	defaultNamespace        = "/registry"
	defaultDialTimeout      = 3 * time.Second
	defaultTTL              = 30 * time.Second
	defaultKeepAliveInterval = 10 * time.Second
)

type options struct {
	endpoints         []string
	dialTimeout       time.Duration
	username          string
	password          string
	tlsConfig         *tls.Config
	namespace         string
	ttl               time.Duration
	keepAliveInterval time.Duration
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		endpoints:         []string{"localhost:2379"},
		dialTimeout:       defaultDialTimeout,
		namespace:         defaultNamespace,
		ttl:               defaultTTL,
		keepAliveInterval: defaultKeepAliveInterval,
	}
}

// Endpoints 设置 etcd 集群地址
func Endpoints(endpoints ...string) Option {
	return func(o *options) {
		if len(endpoints) > 0 {
			o.endpoints = endpoints
		}
	}
}

// DialTimeout 设置连接超时时间，默认 3 秒
func DialTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.dialTimeout = d
		}
	}
}

// Auth 设置用户名和密码认证
func Auth(username, password string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
	}
}

// TLS 设置 TLS 配置
func TLS(cfg *tls.Config) Option {
	return func(o *options) {
		o.tlsConfig = cfg
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
