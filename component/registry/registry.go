package registry

import "context"

// Registry 服务注册与发现接口
type Registry interface {
	// ID 返回实现标识符
	ID() string
	// Register 注册服务
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister 注销服务
	Deregister(ctx context.Context, service *ServiceInstance) error
	// GetService 根据服务名返回服务实例列表
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch 根据服务名创建监听器
	Watch(ctx context.Context, serviceName string) (Watcher, error)
	// Close 关闭,释放资源
	Close() error
}

// Watcher 服务监听器接口
type Watcher interface {
	// Next 在以下两种情况返回服务列表:
	// 1.首次监听且服务实例列表不为空
	// 2.发现任何服务实例变更
	// 如果以上条件都不满足，将阻塞直到上下文超时或取消
	Next() ([]*ServiceInstance, error)
	// Stop 停止监听器
	Stop() error
}
