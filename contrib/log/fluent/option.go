package fluent

import "time"

// Option 是日志配置函数
type Option func(*options)

// options 配置项
type options struct {
	timeout            time.Duration
	writeTimeout       time.Duration
	bufferLimit        int
	retryWait          int
	maxRetry           int
	maxRetryWait       int
	tagPrefix          string
	async              bool
	forceStopAsyncSend bool
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithWriteTimeout 设置写入超时
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = writeTimeout
	}
}

// WithBufferLimit 设置缓冲上限
func WithBufferLimit(bufferLimit int) Option {
	return func(opts *options) {
		opts.bufferLimit = bufferLimit
	}
}

// WithRetryWait 设置重试间隔
func WithRetryWait(retryWait int) Option {
	return func(opts *options) {
		opts.retryWait = retryWait
	}
}

// WithMaxRetry 设置最大重试次数
func WithMaxRetry(maxRetry int) Option {
	return func(opts *options) {
		opts.maxRetry = maxRetry
	}
}

// WithMaxRetryWait 设置最大重试等待
func WithMaxRetryWait(maxRetryWait int) Option {
	return func(opts *options) {
		opts.maxRetryWait = maxRetryWait
	}
}

// WithTagPrefix 设置标签前缀
func WithTagPrefix(tagPrefix string) Option {
	return func(opts *options) {
		opts.tagPrefix = tagPrefix
	}
}

// WithAsync 设置异步发送
func WithAsync(async bool) Option {
	return func(opts *options) {
		opts.async = async
	}
}

// WithForceStopAsyncSend 设置强制停止异步发送
func WithForceStopAsyncSend(forceStopAsyncSend bool) Option {
	return func(opts *options) {
		opts.forceStopAsyncSend = forceStopAsyncSend
	}
}
