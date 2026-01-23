package eureka

import (
	"context"
	"time"
)

type Option func(o *Registry)

// WithContext 设置注册中心上下文
func WithContext(ctx context.Context) Option {
	return func(o *Registry) { o.ctx = ctx }
}

// WithHeartbeat 设置心跳间隔
func WithHeartbeat(interval time.Duration) Option {
	return func(o *Registry) { o.heartbeatInterval = interval }
}

// WithRefresh 设置刷新间隔
func WithRefresh(interval time.Duration) Option {
	return func(o *Registry) { o.refreshInterval = interval }
}

// WithEurekaPath 设置 Eureka 接口路径
func WithEurekaPath(path string) Option {
	return func(o *Registry) { o.eurekaPath = path }
}
