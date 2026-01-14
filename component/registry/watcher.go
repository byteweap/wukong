package registry

// Watcher 服务监听器接口
type Watcher interface {
	// ID 返回实现标识符
	ID() string
	// Next 在以下两种情况返回服务列表:
	// 1.首次监听且服务实例列表不为空
	// 2.发现任何服务实例变更
	// 如果以上条件都不满足，将阻塞直到上下文超时或取消
	Next() ([]*ServiceInstance, error)
	// Stop 关闭监听器
	Stop() error
}
