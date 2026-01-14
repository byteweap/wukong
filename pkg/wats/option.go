package wats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Option 配置选项接口
type Option interface {
	apply(*options)
}

// options 客户端配置选项
type options struct {
	natsOpts []nats.Option
	jsOpts   []jetstream.JetStreamOpt
}

// defaultOptions 返回默认配置
func defaultOptions() *options {
	return &options{
		natsOpts: []nats.Option{},
		jsOpts:   []jetstream.JetStreamOpt{},
	}
}

// natsOptionsWrapper NATS 选项包装器
type natsOptionsWrapper []nats.Option

func (w natsOptionsWrapper) apply(o *options) {
	o.natsOpts = append(o.natsOpts, w...)
}

// WithNats 设置 NATS 连接选项
func WithNats(opts ...nats.Option) Option {
	return natsOptionsWrapper(opts)
}

// jsOptionsWrapper JetStream 选项包装器
type jsOptionsWrapper []jetstream.JetStreamOpt

func (w jsOptionsWrapper) apply(o *options) {
	o.jsOpts = append(o.jsOpts, w...)
}

// WithJetStream 设置 JetStream 选项
func WithJetStream(opts ...jetstream.JetStreamOpt) Option {
	return jsOptionsWrapper(opts)
}
