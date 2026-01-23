package polaris

import "time"

type options struct {
	// 必填，polaris 命名空间
	Namespace string

	// 必填，服务访问 token
	ServiceToken string

	// 可选，协议类型，默认 nil，表示使用服务配置里的协议
	Protocol *string

	// 服务权重，默认 100，范围 0 到 10000
	Weight int

	// 服务优先级，默认 0，值越小优先级越低
	Priority int

	// 是否健康，默认 true
	Healthy bool

	// 是否启用心跳，polaris 本身不提供，默认 true
	Heartbeat bool

	// 是否隔离，默认 false
	Isolate bool

	// TTL 超时，节点使用心跳上报时必填，未设置会返回 ErrorCode-400141
	TTL int

	// 单次查询超时，可选，默认使用全局配置
	// 总超时约为 (1+RetryCount) * Timeout
	Timeout time.Duration

	// 重试次数，可选，默认使用全局配置
	RetryCount int
}

// Option 是 polaris 配置项
type Option func(o *options)

// WithNamespace 设置命名空间
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken 设置服务访问 token
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol 设置协议类型
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithWeight 设置权重
func WithWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy 设置健康状态
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate 设置是否隔离
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithTTL 设置 TTL
func WithTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount 设置重试次数
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

// WithHeartbeat 设置是否启用心跳
func WithHeartbeat(heartbeat bool) Option {
	return func(o *options) { o.Heartbeat = heartbeat }
}
