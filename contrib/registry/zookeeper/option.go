package zookeeper

// Option 是 zookeeper 注册中心配置项
type Option func(o *options)

type options struct {
	namespace string
	user      string
	password  string
}

// WithRootPath 设置注册根路径
func WithRootPath(path string) Option {
	return func(o *options) { o.namespace = path }
}

// WithDigestACL 设置鉴权账号密码
func WithDigestACL(user string, password string) Option {
	return func(o *options) {
		o.user = user
		o.password = password
	}
}
