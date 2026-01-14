package nacos

import (
	"time"
)

const (
	// 默认值
	defaultNamespace      = "public"
	defaultGroup          = "DEFAULT_GROUP"
	defaultClusterName    = "DEFAULT"
	defaultDialTimeout    = 3 * time.Second
	defaultBeatInterval   = 5 * time.Second
	defaultLogLevel       = "info"
	defaultCacheDir       = "/tmp/nacos/cache"
	defaultLogDir         = "/tmp/nacos/log"
	defaultServerPort     = 8848
)

type options struct {
	serverAddrs    []string
	namespace      string
	group          string
	clusterName    string
	dialTimeout    time.Duration
	beatInterval   time.Duration
	username       string
	password       string
	logLevel       string
	cacheDir       string
	logDir         string
	notLoadCacheAtStart bool
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		serverAddrs:         []string{"127.0.0.1:8848"},
		namespace:           defaultNamespace,
		group:               defaultGroup,
		clusterName:         defaultClusterName,
		dialTimeout:         defaultDialTimeout,
		beatInterval:        defaultBeatInterval,
		logLevel:            defaultLogLevel,
		cacheDir:            defaultCacheDir,
		logDir:              defaultLogDir,
		notLoadCacheAtStart: true,
	}
}

// ServerAddrs 设置 Nacos 服务器地址列表，格式: "ip:port" 或 "ip"
func ServerAddrs(addrs ...string) Option {
	return func(o *options) {
		if len(addrs) > 0 {
			o.serverAddrs = addrs
		}
	}
}

// Namespace 设置命名空间，默认 "public"
func Namespace(ns string) Option {
	return func(o *options) {
		if ns != "" {
			o.namespace = ns
		}
	}
}

// Group 设置服务分组，默认 "DEFAULT_GROUP"
func Group(group string) Option {
	return func(o *options) {
		if group != "" {
			o.group = group
		}
	}
}

// ClusterName 设置集群名称，默认 "DEFAULT"
func ClusterName(cluster string) Option {
	return func(o *options) {
		if cluster != "" {
			o.clusterName = cluster
		}
	}
}

// DialTimeout 设置连接超时时间，默认 3 秒
func DialTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.dialTimeout = d
		}
	}
}

// BeatInterval 设置心跳间隔，默认 5 秒
func BeatInterval(interval time.Duration) Option {
	return func(o *options) {
		if interval > 0 {
			o.beatInterval = interval
		}
	}
}

// Auth 设置用户名和密码认证
func Auth(username, password string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
	}
}

// LogLevel 设置日志级别，默认 "info"
func LogLevel(level string) Option {
	return func(o *options) {
		if level != "" {
			o.logLevel = level
		}
	}
}

// CacheDir 设置缓存目录，默认 "/tmp/nacos/cache"
func CacheDir(dir string) Option {
	return func(o *options) {
		if dir != "" {
			o.cacheDir = dir
		}
	}
}

// LogDir 设置日志目录，默认 "/tmp/nacos/log"
func LogDir(dir string) Option {
	return func(o *options) {
		if dir != "" {
			o.logDir = dir
		}
	}
}

// NotLoadCacheAtStart 设置启动时不加载缓存，默认 true
func NotLoadCacheAtStart(notLoad bool) Option {
	return func(o *options) {
		o.notLoadCacheAtStart = notLoad
	}
}
