package locator

import "context"

// Locator 跟踪玩家会话在各节点间的位置
type Locator interface {

	// ID 返回定位器实现标识符
	ID() string

	// Node 返回用户所在某服务节点
	Node(ctx context.Context, uid int64, service string) (string, error)

	// Bind 绑定用户到某服务某节点
	Bind(ctx context.Context, uid int64, service, node string) error

	// UnBind 解绑用户的某服务某节点
	UnBind(ctx context.Context, uid int64, service, node string) error

	// Close 关闭定位器
	Close() error
}
