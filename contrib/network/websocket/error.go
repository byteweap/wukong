package websocket

import "errors"

// WebSocket错误定义
// 定义了WebSocket实现中可能返回的各种错误类型
var (
	// ErrInvalidOpCode 无效的操作码错误
	// 当收到的WebSocket帧的操作码不符合规范时返回
	ErrInvalidOpCode = errors.New("invalid opcode")

	// ErrInvalidState 无效的连接状态错误
	// 当连接处于非法状态时返回
	ErrInvalidState = errors.New("invalid state")

	// ErrHubClosed Hub已关闭错误
	// 当尝试向已关闭的Hub分配新连接时返回
	ErrHubClosed = errors.New("hub closed")

	// ErrWriteQueueFull 写入队列已满错误
	// 当消息写入队列达到容量上限时返回，通常是因为发送速度跟不上接收速度
	ErrWriteQueueFull = errors.New("write queue full")

	// ErrMaxConns 连接数达到上限错误
	// 当尝试建立新连接但已达到最大连接限制时返回
	ErrMaxConns = errors.New("max connections reached")
)
