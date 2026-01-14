package nacos

import (
	"context"
	"fmt"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/byteweap/wukong/component/registry"
)

// ID Nacos 发现器实现标识符
const DiscoveryID = "nacos(discovery)"

// Discovery 使用 Nacos 实现服务发现
type Discovery struct {
	namingClient naming_client.INamingClient
	opts         *options
}

var _ registry.Discovery = (*Discovery)(nil)

// NewDiscovery 创建 Nacos 发现器并立即连接
func NewDiscovery(opts ...Option) (*Discovery, error) {
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

	return &Discovery{
		namingClient: namingClient,
		opts:         o,
	}, nil
}

// NewDiscoveryWith 使用现有的 Nacos 客户端创建发现器（调用者负责其生命周期）
func NewDiscoveryWith(namingClient naming_client.INamingClient, opts ...Option) *Discovery {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Discovery{
		namingClient: namingClient,
		opts:         o,
	}
}

// ID 返回实现标识符
func (d *Discovery) ID() string {
	return DiscoveryID
}

// GetService 根据服务名返回服务实例列表
func (d *Discovery) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// 查询参数
	param := vo.SelectAllInstancesParam{
		ServiceName: serviceName,
		GroupName:   d.opts.group,
		Clusters:    []string{d.opts.clusterName},
	}
	// 获取服务实例
	instances, err := d.namingClient.SelectAllInstances(param)
	if err != nil {
		return nil, fmt.Errorf("failed to get services from nacos: %w", err)
	}
	// 转换为 ServiceInstance
	result := make([]*registry.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		si := convertToServiceInstance(instance, serviceName)
		if si != nil {
			result = append(result, si)
		}
	}

	return result, nil
}

// Watch 根据服务名创建监听器
func (d *Discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	return newWatcher(ctx, d.namingClient, serviceName, d.opts)
}

// Close 关闭发现器并释放资源
func (d *Discovery) Close() error {
	// Nacos 客户端没有显式的 Close 方法
	return nil
}

// convertToServiceInstance 将 Nacos Instance 转换为 ServiceInstance
func convertToServiceInstance(instance model.Instance, serviceName string) *registry.ServiceInstance {
	ip := instance.Ip
	port := instance.Port
	metadata := instance.Metadata
	instanceId := instance.InstanceId

	// 构建 endpoints
	endpoints := []string{}
	if metadata != nil {
		// 优先使用 metadata 中的 endpoints
		if eps, ok := metadata["endpoints"]; ok && eps != "" {
			endpoints = strings.Split(eps, ",")
		}
	}

	// 如果没有 endpoints，从 IP:Port 构建
	if len(endpoints) == 0 {
		// 尝试从 metadata 获取协议
		protocol := "http"
		if p, ok := metadata["protocol"]; ok {
			protocol = p
		}
		endpoint := fmt.Sprintf("%s://%s:%d", protocol, ip, port)
		endpoints = []string{endpoint}
	}

	// 构建 ServiceInstance
	si := &registry.ServiceInstance{
		ID:        instanceId,
		Name:      serviceName,
		Endpoints: endpoints,
		Metadata:  make(map[string]string),
	}

	// 复制 metadata
	if metadata != nil {
		for k, v := range metadata {
			// 跳过内部使用的字段
			if k != "endpoints" {
				si.Metadata[k] = v
			}
		}
		// 提取版本
		if version, ok := metadata["version"]; ok {
			si.Version = version
		}
	}

	// 如果没有 InstanceId，使用 IP:Port 作为 ID
	if si.ID == "" {
		si.ID = fmt.Sprintf("%s:%d", ip, port)
	}

	return si
}
