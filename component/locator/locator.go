package locator

import "context"

// Locator 跟踪玩家会话在各节点间的位置
type Locator interface {

	// ID 返回定位器实现标识符
	ID() string

	// Gate 返回用户ID对应的网关节点
	Gate(ctx context.Context, uid int64) (string, error)

	// BindGate 绑定用户ID到网关节点
	BindGate(ctx context.Context, uid int64, node string) error

	// UnBindGate 解绑用户ID的网关节点
	UnBindGate(ctx context.Context, uid int64, node string) error

	// Game 返回用户ID对应的游戏节点
	Game(ctx context.Context, uid int64) (string, error)

	// BindGame 绑定用户ID到游戏节点
	BindGame(ctx context.Context, uid int64, node string) error

	// UnBindGame 解绑用户ID的游戏节点
	UnBindGame(ctx context.Context, uid int64, node string) error

	// Close 关闭定位器
	Close() error
}
