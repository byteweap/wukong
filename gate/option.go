package gate

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	// Addr is the address to listen on.
	Addr string
	// Pattern is the pattern to listen on.
	Pattern string
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

	// LocatorKeyFormat is the key format for the locator.
	LocatorKeyFormat string // e.g. "gate:%d"
	// LocatorGateFieldName is the field name for the gate in the locator.
	LocatorGateFieldName string // e.g. "gate"
	// LocatorGameFieldName is the field name for the game in the locator.
	LocatorGameFieldName string // e.g. "game"

	// RedisOptions is the Redis options.
	RedisOptions *redis.UniversalOptions
}

type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Addr:                 "0.0.0.0:8000",
		Pattern:              "/",
		MaxConnections:       10000,
		MaxMessageSize:       4 * 1024, // 4KB
		ReadTimeout:          0,
		WriteTimeout:         0,
		WriteQueueSize:       0,
		LocatorKeyFormat:     "gate:%d",
		LocatorGateFieldName: "gate",
		LocatorGameFieldName: "game",
		RedisOptions: &redis.UniversalOptions{
			Addrs:      []string{"localhost:6379"},
			Username:   "",
			Password:   "",
			DB:         0,
			ClientName: "wukong-gate",
		},
	}
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithPattern(pattern string) Option {
	return func(o *Options) {
		if o.Pattern != "" {
			o.Pattern = pattern
		}
	}
}

func WithMaxConnections(maxConnections int) Option {
	return func(o *Options) {
		if o.MaxConnections > 0 {
			o.MaxConnections = maxConnections
		}
	}
}
