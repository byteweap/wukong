package tencent

// Option 配置函数
type Option func(cls *options)

// options 配置项
type options struct {
	topicID      string
	accessKey    string
	accessSecret string
	endpoint     string
}

// defaultOptions 默认配置
func defaultOptions() *options {
	return &options{}
}

// WithEndpoint 设置服务地址
func WithEndpoint(endpoint string) Option {
	return func(cls *options) {
		cls.endpoint = endpoint
	}
}

// WithTopicID 设置主题 ID
func WithTopicID(topicID string) Option {
	return func(cls *options) {
		cls.topicID = topicID
	}
}

// WithAccessKey 设置访问键
func WithAccessKey(ak string) Option {
	return func(cls *options) {
		cls.accessKey = ak
	}
}

// WithAccessSecret 设置访问密钥
func WithAccessSecret(as string) Option {
	return func(cls *options) {
		cls.accessSecret = as
	}
}
