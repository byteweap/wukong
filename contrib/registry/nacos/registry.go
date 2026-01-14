package nacos

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/byteweap/wukong/component/registry"
)

// ID Nacos 注册器实现标识符
const RegistryID = "nacos(registry)"

// Registry 使用 Nacos 实现服务注册
type Registry struct {
	namingClient  naming_client.INamingClient
	opts          *options
	beatCancelMap sync.Map // map[string]context.CancelFunc 存储实例ID对应的取消函数
}

var _ registry.Registrar = (*Registry)(nil)

// NewRegistry 创建 Nacos 注册器并立即连接
func NewRegistry(opts ...Option) (*Registry, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	// 解析服务器地址
	serverConfigs := make([]constant.ServerConfig, 0, len(o.serverAddrs))
	for _, addr := range o.serverAddrs {
		ip, port := parseAddr(addr)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr:      ip,
			Port:        port,
			ContextPath: "/nacos",
		})
	}

	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         o.namespace,
		TimeoutMs:           uint64(o.dialTimeout.Milliseconds()),
		NotLoadCacheAtStart: o.notLoadCacheAtStart,
		LogDir:              o.logDir,
		CacheDir:            o.cacheDir,
		LogLevel:            o.logLevel,
	}

	if o.username != "" {
		clientConfig.Username = o.username
		clientConfig.Password = o.password
	}

	// 创建命名客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos naming client: %w", err)
	}

	return &Registry{
		namingClient: namingClient,
		opts:         o,
	}, nil
}

// NewRegistryWith 使用现有的 Nacos 客户端创建注册器（调用者负责其生命周期）
func NewRegistryWith(namingClient naming_client.INamingClient, opts ...Option) *Registry {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Registry{
		namingClient: namingClient,
		opts:         o,
	}
}

// ID 返回实现标识符
func (r *Registry) ID() string {
	return RegistryID
}

// Register 注册服务实例
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	if service == nil {
		return fmt.Errorf("service instance is nil")
	}
	if service.ID == "" || service.Name == "" {
		return fmt.Errorf("service ID and Name are required")
	}
	if len(service.Endpoints) == 0 {
		return fmt.Errorf("service endpoints are required")
	}

	// 解析第一个 endpoint 获取 IP 和端口
	ip, port, err := parseEndpoint(service.Endpoints[0])
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	// 构建元数据
	metadata := make(map[string]string)
	if service.Metadata != nil {
		for k, v := range service.Metadata {
			metadata[k] = v
		}
	}
	// 添加版本信息
	if service.Version != "" {
		metadata["version"] = service.Version
	}
	// 添加所有 endpoints
	metadata["endpoints"] = strings.Join(service.Endpoints, ",")

	// 注册参数
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: service.Name,
		Weight:      service.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, // 临时实例，支持自动下线
		Metadata:    metadata,
		GroupName:   r.opts.group,
		ClusterName: r.opts.clusterName,
	}

	// 注册服务
	success, err := r.namingClient.RegisterInstance(param)
	if err != nil {
		return fmt.Errorf("failed to register instance: %w", err)
	}
	if !success {
		return fmt.Errorf("register instance failed: unknown error")
	}

	// 注意：对于 Ephemeral=true 的实例，Nacos SDK 会自动处理心跳
	// 我们只需要存储取消函数用于资源清理（虽然当前不需要，但保留结构以便将来扩展）
	_, cancel := context.WithCancel(context.Background())
	r.beatCancelMap.Store(service.ID, cancel)

	return nil
}

// Deregister 注销服务实例
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	if service == nil {
		return fmt.Errorf("service instance is nil")
	}
	if service.ID == "" || service.Name == "" {
		return fmt.Errorf("service ID and Name are required")
	}

	// 停止心跳
	if cancel, ok := r.beatCancelMap.LoadAndDelete(service.ID); ok {
		cancel.(context.CancelFunc)()
	}

	// 解析第一个 endpoint 获取 IP 和端口
	ip, port, err := parseEndpoint(service.Endpoints[0])
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	// 注销参数
	param := vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: service.Name,
		GroupName:   r.opts.group,
		Cluster:     r.opts.clusterName,
		Ephemeral:   true,
	}

	// 注销服务
	success, err := r.namingClient.DeregisterInstance(param)
	if err != nil {
		return fmt.Errorf("failed to deregister instance: %w", err)
	}
	if !success {
		return fmt.Errorf("deregister instance failed: unknown error")
	}

	return nil
}

// Close 关闭注册器并释放资源
func (r *Registry) Close() error {
	// 停止所有心跳
	r.beatCancelMap.Range(func(key, value interface{}) bool {
		cancel := value.(context.CancelFunc)
		cancel()
		return true
	})

	// Nacos 客户端没有显式的 Close 方法，这里只清理心跳
	return nil
}

// 注意：Nacos SDK 对于 Ephemeral=true 的实例会自动处理心跳
// 这里保留 keepAlive 和 sendBeat 的占位符，以便将来扩展自定义心跳逻辑

// parseAddr 解析地址，格式: "ip:port" 或 "ip"
func parseAddr(addr string) (string, uint64) {
	parts := strings.Split(addr, ":")
	ip := parts[0]
	port := defaultServerPort
	if len(parts) > 1 {
		if p, err := strconv.Atoi(parts[1]); err == nil {
			port = p
		}
	}
	return ip, uint64(port)
}

// parseEndpoint 解析 endpoint，格式: "http://ip:port" 或 "grpc://ip:port" 或 "ip:port"
func parseEndpoint(endpoint string) (string, int, error) {
	// 尝试解析为 URL
	u, err := url.Parse(endpoint)
	if err == nil && u.Host != "" {
		host, portStr, err := net.SplitHostPort(u.Host)
		if err != nil {
			// 如果没有端口，使用默认端口
			host = u.Host
			portStr = "80"
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return "", 0, fmt.Errorf("invalid port: %s", portStr)
		}
		return host, port, nil
	}

	// 尝试直接解析为 "ip:port"
	host, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "", 0, fmt.Errorf("invalid endpoint format: %s", endpoint)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port: %s", portStr)
	}
	return host, port, nil
}
