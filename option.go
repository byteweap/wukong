package wukong

import (
	"context"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/server"
	"github.com/google/uuid"
)

// Option 为应用选项
type Option func(o *options)

// options 为应用配置项集合
type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	logger          log.Logger
	registry        registry.Registry
	registryTimeout time.Duration
	stopTimeout     time.Duration
	servers         []server.Server

	// 启停前后回调
	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// defaultOptions 返回默认配置
func defaultOptions() *options {
	return &options{
		id:              uuid.New().String(),
		name:            "wukong",
		version:         "v1.0.0",
		metadata:        map[string]string{},
		endpoints:       []*url.URL{},
		ctx:             context.Background(),
		sigs:            []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registryTimeout: time.Second * 10,
	}
}

// ID 设置服务 ID
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name 设置服务名
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version 设置服务版本
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata 设置服务元数据
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Endpoint 设置服务端点
func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context 设置服务上下文
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Logger 设置服务日志器
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Signal 设置退出信号
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Registry 设置服务注册中心
func Registry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// RegistryTimeout 设置注册超时
func RegistryTimeout(t time.Duration) Option {
	return func(o *options) { o.registryTimeout = t }
}

// StopTimeout 设置停止超时
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// 启停前后回调

// BeforeStart 添加启动前回调
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// BeforeStop 添加停止前回调
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStart 添加启动后回调
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// AfterStop 添加停止后回调
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
