package nacos

import (
	"time"
)

const (
	// 默认值
	defaultNamespace   = "public"
	defaultGroup       = "DEFAULT_GROUP"
	defaultClusterName = "DEFAULT"
	defaultDialTimeout = 3 * time.Second
	defaultLogLevel    = "info"
	defaultCacheDir    = "/tmp/nacos/cache"
	defaultLogDir      = "/tmp/nacos/log"
	defaultServerPort  = 8848
)

// NacosConfig Nacos 配置
type NacosConfig struct {
	// Addrs 服务器地址列表，格式: "ip:port" 或 "ip"
	Addrs []string
	// Namespace 命名空间，默认 "public"
	Namespace string
	// DialTimeout 连接超时时间，默认 3 秒
	DialTimeout time.Duration
	// Username 用户名（可选）
	Username string
	// Password 密码（可选）
	Password string
	// LogLevel 日志级别，默认 "info"
	LogLevel string
	// CacheDir 缓存目录，默认 "/tmp/nacos/cache"
	CacheDir string
	// LogDir 日志目录，默认 "/tmp/nacos/log"
	LogDir string
	// NotLoadCacheAtStart 启动时不加载缓存，默认 true
	NotLoadCacheAtStart bool
}

// DefaultNacosConfig 返回默认的 Nacos 配置
func DefaultNacosConfig() *NacosConfig {
	return &NacosConfig{
		Addrs:               []string{"127.0.0.1:8848"},
		Namespace:           defaultNamespace,
		DialTimeout:         defaultDialTimeout,
		LogLevel:            defaultLogLevel,
		CacheDir:            defaultCacheDir,
		LogDir:              defaultLogDir,
		NotLoadCacheAtStart: true,
	}
}

func validate(cfg *NacosConfig) {
	if len(cfg.Addrs) == 0 {
		cfg.Addrs = []string{"127.0.0.1:8848"}
	}
	if cfg.Namespace == "" {
		cfg.Namespace = defaultNamespace
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = defaultDialTimeout
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = defaultLogLevel
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = defaultCacheDir
	}
	if cfg.LogDir == "" {
		cfg.LogDir = defaultLogDir
	}
	if cfg.NotLoadCacheAtStart {
		cfg.NotLoadCacheAtStart = true
	}
}

// options 客户端配置选项
type options struct {
	group       string
	clusterName string
	weight      float64 // 服务注册权重, 必须大于0
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		group:       defaultGroup,
		clusterName: defaultClusterName,
		weight:      10,
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

// Weight 设置服务注册权重, 必须大于0, 默认10
func Weight(weight float64) Option {
	return func(o *options) {
		if weight > 0 {
			o.weight = weight
		}
	}
}
