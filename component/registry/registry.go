package registry

import "context"

// Registrar 服务注册接口
type Registrar interface {
	// ID 返回实现标识符
	ID() string
	// Register 注册服务
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister 注销服务
	Deregister(ctx context.Context, service *ServiceInstance) error
	// Close 关闭,释放资源
	Close() error
}
