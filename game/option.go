package game

import (
	"context"

	"github.com/byteweap/wukong/component/log"
	"github.com/google/uuid"
)

type (
	application struct {
		id       string
		name     string
		version  string
		metadata map[string]string
	}
)
type options struct {
	ctx    context.Context
	app    application
	logger log.Logger
}

func defaultOptions() *options {
	return &options{
		ctx: context.Background(),
		app: application{
			id:       uuid.New().String(),
			name:     "wukong-game",
			version:  "v1.0.0",
			metadata: make(map[string]string),
		},
	}
}

type Option func(*options)

// Context 设置上下文, 默认: context.Background()
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// ID 设置应用ID, 默认: uuid()
func ID(id string) Option {
	return func(o *options) { o.app.id = id }
}

// Name 设置应用名称, 默认: "application"
func Name(name string) Option {
	return func(o *options) { o.app.name = name }
}

// Version 设置应用版本, 默认: "v1.0.0"
func Version(version string) Option {
	return func(o *options) { o.app.version = version }
}

// Metadata 设置应用元数据
func Metadata(metadata map[string]string) Option {
	return func(o *options) { o.app.metadata = metadata }
}

// Logger 设置日志记录器
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}
