package etcd

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/byteweap/wukong/component/registry"
	"github.com/byteweap/wukong/pkg/kcodec"
)

// ID etcd 注册器实现标识符
const RegistryID = "etcd"

// Registry 使用 etcd 实现服务注册
type Registry struct {
	client            *clientv3.Client
	namespace         string
	ttl               time.Duration
	keepAliveInterval time.Duration
	leases            sync.Map // map[string]*leaseInfo 存储实例ID对应的租约信息
	ownClient         bool     // 是否由自己管理客户端的生命周期
}

var _ registry.Registry = (*Registry)(nil)

type leaseInfo struct {
	leaseID clientv3.LeaseID
	cancel  context.CancelFunc
}

// New 创建 etcd 注册器并立即连接
func NewRegistry(opts ...Option) (*Registry, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	cfg := clientv3.Config{
		Endpoints:   o.endpoints,
		DialTimeout: o.dialTimeout,
		Username:    o.username,
		Password:    o.password,
		// TLS:         o.tlsConfig,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &Registry{
		client:            client,
		namespace:         o.namespace,
		ttl:               o.ttl,
		keepAliveInterval: o.keepAliveInterval,
		ownClient:         true, // NewRegistry 创建的实例自己管理客户端生命周期
	}, nil
}

// NewWith 使用现有的 etcd 客户端创建注册器（调用者负责其生命周期）
func NewRegistryWith(client *clientv3.Client, opts ...Option) *Registry {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Registry{
		client:            client,
		namespace:         o.namespace,
		ttl:               o.ttl,
		keepAliveInterval: o.keepAliveInterval,
		ownClient:         false, // NewRegistryWith 创建的实例由上层管理客户端生命周期
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

	// 序列化服务实例
	data, err := kcodec.Invoke("json")
	if err != nil {
		return fmt.Errorf("failed to get json codec: %w", err)
	}
	value, err := data.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service instance: %w", err)
	}

	// 构建 key
	key := r.buildKey(service.Name, service.ID)

	// 创建租约
	leaseResp, err := r.client.Grant(ctx, int64(r.ttl.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to create lease: %w", err)
	}

	// 存储 key-value，并绑定租约
	_, err = r.client.Put(ctx, key, string(value), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to put key-value: %w", err)
	}

	// 启动续租 goroutine
	keepAliveCtx, cancel := context.WithCancel(context.Background())
	keepAliveResp, err := r.client.KeepAlive(keepAliveCtx, leaseResp.ID)
	if err != nil {
		cancel()
		// KeepAlive 失败，撤销租约并删除 key
		_, _ = r.client.Revoke(context.Background(), leaseResp.ID)
		return fmt.Errorf("failed to start keepalive: %w", err)
	}

	// 存储租约信息
	r.leases.Store(service.ID, &leaseInfo{
		leaseID: leaseResp.ID,
		cancel:  cancel,
	})

	// 启动 goroutine 处理续租响应（避免 channel 阻塞）
	go func() {
		for {
			select {
			case <-keepAliveCtx.Done():
				return
			case resp, ok := <-keepAliveResp:
				if !ok {
					// channel 关闭，可能是 etcd 连接断开
					return
				}
				if resp == nil {
					// 续租失败
					return
				}
				// 续租成功，继续等待下一次
			}
		}
	}()

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

	// 获取租约信息
	info, ok := r.leases.LoadAndDelete(service.ID)
	if !ok {
		// 如果没有租约信息，直接删除 key
		key := r.buildKey(service.Name, service.ID)
		_, err := r.client.Delete(ctx, key)
		return err
	}

	leaseInfo := info.(*leaseInfo)

	// 停止续租
	leaseInfo.cancel()

	// 撤销租约（会自动删除关联的 key）
	_, err := r.client.Revoke(ctx, leaseInfo.leaseID)
	if err != nil {
		// 如果撤销失败，尝试直接删除 key
		key := r.buildKey(service.Name, service.ID)
		_, _ = r.client.Delete(ctx, key)
		return fmt.Errorf("failed to revoke lease: %w", err)
	}

	return nil
}

// GetService 根据服务名返回服务实例列表
func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// 构建 key 前缀
	prefix := r.buildKeyPrefix(serviceName)

	// 获取所有匹配的 key-value
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
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
func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// 构建 key 前缀
	prefix := r.buildKeyPrefix(serviceName)

	// 创建 watcher
	return newWatcher(ctx, r.client, prefix, r.namespace)
}

// Close 关闭注册器并释放资源
func (r *Registry) Close() error {
	// 停止所有续租
	r.leases.Range(func(key, value interface{}) bool {
		info := value.(*leaseInfo)
		info.cancel()
		return true
	})

	// 只有自己管理的客户端才关闭
	if r.ownClient && r.client != nil {
		return r.client.Close()
	}
	return nil
}

// buildKey 构建 etcd key
func (r *Registry) buildKey(serviceName, instanceID string) string {
	return path.Join(r.namespace, serviceName, instanceID)
}

// buildKeyPrefix 构建 etcd key 前缀
func (r *Registry) buildKeyPrefix(serviceName string) string {
	return path.Join(r.namespace, serviceName) + "/"
}
