package gate

import "time"

type Options struct {
	// Addr is the address to listen on.
	Addr string
	// Port is the port to listen on.
	Port int
	// MaxMessageSize is the maximum message size.
	MaxMessageSize int64
	// MaxConnections is the maximum number of connections.
	MaxConnections int
	// ReadTimeout is the read timeout.
	ReadTimeout time.Duration
	// WriteTimeout is the write timeout.
	WriteTimeout time.Duration
	// WriteQueueSize is the write queue size.
	WriteQueueSize int
}

type Option func(*Options)

var defaultOptions = Options{
	Addr:           "0.0.0.0",
	Port:           8080,
	MaxConnections: 1024,
}

func newOptions(opts ...Option) Options {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func WithMaxConnections(maxConnections int) Option {
	return func(o *Options) {
		o.MaxConnections = maxConnections
	}
}
