package gate

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	// NetworkOptions is the options for the network.
	NetworkOptions struct {
		Addr           string
		Pattern        string
		MaxMessageSize int64
		MaxConnections int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		WriteQueueSize int
	}

	// LocatorOptions is the options for the locator.
	LocatorOptions struct {
		KeyFormat     string
		GateFieldName string
		GameFieldName string
	}

	// Options is the options for the gate.
	Options struct {
		NetworkOptions NetworkOptions
		LocatorOptions LocatorOptions
		RedisOptions   redis.UniversalOptions
	}
)

type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		NetworkOptions: NetworkOptions{
			Addr:           "0.0.0.0:8000",
			Pattern:        "/",
			MaxConnections: 10000,
			MaxMessageSize: 4 * 1024, // 4KB
			ReadTimeout:    0,
			WriteTimeout:   0,
			WriteQueueSize: 0,
		},
		LocatorOptions: LocatorOptions{
			KeyFormat:     "gate:%d",
			GateFieldName: "gate",
			GameFieldName: "game",
		},
		RedisOptions: redis.UniversalOptions{
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
		o.NetworkOptions.Addr = addr
	}
}

func WithPattern(pattern string) Option {
	return func(o *Options) {
		if o.NetworkOptions.Pattern != "" {
			o.NetworkOptions.Pattern = pattern
		}
	}
}

func WithMaxConnections(maxConnections int) Option {
	return func(o *Options) {
		if o.NetworkOptions.MaxConnections > 0 {
			o.NetworkOptions.MaxConnections = maxConnections
		}
	}
}

func WithLocator(keyFormat, gateFieldName, gameFieldName string) Option {
	return func(o *Options) {
		o.LocatorOptions.KeyFormat = keyFormat
		o.LocatorOptions.GateFieldName = gateFieldName
		o.LocatorOptions.GameFieldName = gameFieldName
	}
}
