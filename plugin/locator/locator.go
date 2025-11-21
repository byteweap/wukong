package locator

// Locator player locator
type Locator interface {

	// ID 标识，区分不同的实现
	ID() string

	// Gate 获取玩家当前Gate节点
	Gate(uid int64) (string, error)
	// BindGate 绑定Gate节点
	BindGate(uid int64, node string) error
	// UnBindGate 解绑Gate节点
	UnBindGate(uid int64, node string) error

	// Game 获取玩家当前Game节点
	Game(uid int64) (string, error)
	// BindGame 绑定Game节点
	BindGame(uid int64, node string) error
	// UnBindGame 解绑Game节点
	UnBindGame(uid int64, node string) error
}
