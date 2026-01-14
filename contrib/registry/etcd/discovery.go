package etcd

import (
	"context"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/pkg/kcodec"
)

// ID etcd 发现器实现标识符
const DiscoveryID = "etcd(discovery)"

// Discovery 使用 etcd 实现服务发现
type Discovery struct {
	client    *clientv3.Client
	namespace string
}

var _ registry.Discovery = (*Discovery)(nil)

// NewDiscovery 创建 etcd 发现器并立即连接
func NewDiscovery(opts ...Option) (*Discovery, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	cfg := clientv3.Config{
		Endpoints:   o.endpoints,
		DialTimeout: o.dialTimeout,
		Username:    o.username,
		Password:    o.password,
		TLS:         o.tlsConfig,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &Discovery{
		client:    client,
		namespace: o.namespace,
	}, nil
}

// NewDiscoveryWith 使用现有的 etcd 客户端创建发现器（调用者负责其生命周期）
func NewDiscoveryWith(client *clientv3.Client, opts ...Option) *Discovery {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Discovery{
		client:    client,
		namespace: o.namespace,
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

	// 构建 key 前缀
	prefix := d.buildKeyPrefix(serviceName)

	// 获取所有匹配的 key-value
	resp, err := d.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get services from etcd: %w", err)
	}

	// 解析服务实例
	instances := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	codec, err := kcodec.Invoke("json")
	if err != nil {
		return nil, fmt.Errorf("failed to get json codec: %w", err)
	}

	for _, kv := range resp.Kvs {
		var instance registry.ServiceInstance
		if err := codec.Unmarshal(kv.Value, &instance); err != nil {
			// 跳过无法解析的实例
			continue
		}
		instances = append(instances, &instance)
	}

	return instances, nil
}

// Watch 根据服务名创建监听器
func (d *Discovery) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// 构建 key 前缀
	prefix := d.buildKeyPrefix(serviceName)

	// 创建 watcher
	return newWatcher(ctx, d.client, prefix, d.namespace)
}

// Close 关闭发现器并释放资源
func (d *Discovery) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

// buildKeyPrefix 构建 etcd key 前缀
func (d *Discovery) buildKeyPrefix(serviceName string) string {
	return path.Join(d.namespace, serviceName) + "/"
}
