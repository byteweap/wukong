package aliyun

// Option 配置函数
type Option func(alc *options)

// options 配置项
type options struct {
	accessKey    string
	accessSecret string
	endpoint     string
	project      string
	logstore     string
}

// defaultOptions 默认配置
func defaultOptions() *options {
	return &options{
		project:  "projectName",
		logstore: "app",
	}
}

// WithEndpoint 设置服务地址
func WithEndpoint(endpoint string) Option {
	return func(alc *options) {
		alc.endpoint = endpoint
	}
}

// WithProject 设置项目名
func WithProject(project string) Option {
	return func(alc *options) {
		alc.project = project
	}
}

// WithLogstore 设置日志库
func WithLogstore(logstore string) Option {
	return func(alc *options) {
		alc.logstore = logstore
	}
}

// WithAccessKey 设置访问键
func WithAccessKey(ak string) Option {
	return func(alc *options) {
		alc.accessKey = ak
	}
}

// WithAccessSecret 设置访问密钥
func WithAccessSecret(as string) Option {
	return func(alc *options) {
		alc.accessSecret = as
	}
}
