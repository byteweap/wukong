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

func defaultOptions() *Options {
	return &Options{
		Addr:           "0.0.0.0",
		Port:           8080,
		MaxConnections: 10000,
		MaxMessageSize: 4 * 1024, // 4KB
		ReadTimeout:    0,
		WriteTimeout:   0,
		WriteQueueSize: 0,
	}
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
