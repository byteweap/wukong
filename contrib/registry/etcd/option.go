package etcd

import (
	"context"
	"time"
)

// Option 是 etcd 注册中心配置项
type Option func(o *options)

type options struct {
	ctx       context.Context
	namespace string
	ttl       time.Duration
	maxRetry  int
}

// Context 设置注册中心上下文
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Namespace 设置注册中心命名空间
func Namespace(ns string) Option {
	return func(o *options) { o.namespace = ns }
}

// RegisterTTL 设置注册 TTL
func RegisterTTL(ttl time.Duration) Option {
	return func(o *options) { o.ttl = ttl }
}

// MaxRetry 设置最大重试次数
func MaxRetry(num int) Option {
	return func(o *options) { o.maxRetry = num }
}
