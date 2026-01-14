package registry

import "context"

// Discovery 服务发现接口
type Discovery interface {
	// ID 返回实现标识符
	ID() string
	// GetService 根据服务名返回内存中的服务实例列表
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch 根据服务名创建监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}
