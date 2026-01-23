package nacos

type options struct {
	prefix  string
	weight  float64
	cluster string
	group   string
	kind    string
}

// Option 是 nacos 配置项
type Option func(o *options)

// WithPrefix 设置前缀路径
func WithPrefix(prefix string) Option {
	return func(o *options) { o.prefix = prefix }
}

// WithWeight 设置权重
func WithWeight(weight float64) Option {
	return func(o *options) { o.weight = weight }
}

// WithCluster 设置集群名
func WithCluster(cluster string) Option {
	return func(o *options) { o.cluster = cluster }
}

// WithGroup 设置分组
func WithGroup(group string) Option {
	return func(o *options) { o.group = group }
}

// WithDefaultKind 设置默认协议类型
func WithDefaultKind(kind string) Option {
	return func(o *options) { o.kind = kind }
}
