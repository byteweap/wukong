package gate

import (
	"context"
	"time"

	"github.com/byteweap/wukong/log"
)

type (
	// ApplicationOptions 应用选项
	ApplicationOptions struct {
		ID       string
		Name     string
		Version  string
		Metadata map[string]string
		Addr     string
	}

	// options 选项
	options struct {
		ctx             context.Context
		logger          log.Logger
		application     ApplicationOptions
		registryTimeout time.Duration
	}
)

type Option func(*options)

func defaultOptions() *options {

	return &options{
		application: ApplicationOptions{
			ID:       "",
			Name:     "gate",
			Version:  "1.0.0",
			Metadata: make(map[string]string),
			Addr:     "0.0.0.0:9000",
		},
	}
}

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}
func ID(id string) Option {
	return func(o *options) {
		o.application.ID = id
	}
}

func Name(name string) Option {
	return func(o *options) {
		o.application.Name = name
	}
}

func Version(version string) Option {
	return func(o *options) {
		o.application.Version = version
	}
}

func Metadata(metadata map[string]string) Option {
	return func(o *options) {
		o.application.Metadata = metadata
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func RegistryTimeout(registryTimeout time.Duration) Option {
	return func(o *options) {
		o.registryTimeout = registryTimeout
	}
}
