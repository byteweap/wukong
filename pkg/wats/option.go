package wats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Option interface {
	apply(*options)
}

type options struct {
	natsOpts []nats.Option
	jsOpts   []jetstream.JetStreamOpt
}

func defaultOptions() *options {
	return &options{
		natsOpts: []nats.Option{},
		jsOpts:   []jetstream.JetStreamOpt{},
	}
}

type natsOptionsWrapper []nats.Option

func (w natsOptionsWrapper) apply(o *options) {
	o.natsOpts = append(o.natsOpts, w...)
}

func WithNats(opts ...nats.Option) Option {
	return natsOptionsWrapper(opts)
}

type jsOptionsWrapper []jetstream.JetStreamOpt

func (w jsOptionsWrapper) apply(o *options) {
	o.jsOpts = append(o.jsOpts, w...)
}

func WithJetStream(opts ...jetstream.JetStreamOpt) Option {
	return jsOptionsWrapper(opts)
}
