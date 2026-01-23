package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
)

// Option 是 consul 注册中心的配置项
type Option func(*Registry)

// WithHealthCheck 设置是否启用健康检查
func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

// WithTimeout 设置获取服务的超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *Registry) {
		o.timeout = timeout
	}
}

// WithDatacenter 设置数据中心
func WithDatacenter(dc Datacenter) Option {
	return func(o *Registry) {
		o.cli.dc = dc
	}
}

// WithHeartbeat 设置是否启用心跳
func WithHeartbeat(enable bool) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.heartbeat = enable
		}
	}
}

// WithServiceResolver 设置 endpoint 解析函数
func WithServiceResolver(fn ServiceResolver) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.resolver = fn
		}
	}
}

// WithHealthCheckInterval 设置健康检查间隔秒数
func WithHealthCheckInterval(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.healthcheckInterval = interval
		}
	}
}

// WithDeregisterCriticalServiceAfter 设置不健康服务自动注销时间，单位秒
func WithDeregisterCriticalServiceAfter(interval int) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.deregisterCriticalServiceAfter = interval
		}
	}
}

// WithServiceCheck 追加服务健康检查
func WithServiceCheck(checks ...*api.AgentServiceCheck) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.serviceChecks = checks
		}
	}
}

// WithTags 设置服务标签
func WithTags(tags []string) Option {
	return func(o *Registry) {
		if o.cli != nil {
			o.cli.tags = tags
		}
	}
}
