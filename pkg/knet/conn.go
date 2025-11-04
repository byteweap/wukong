package knet

import "net"

// Conn 网络连接的核心接口
// 所有网络协议（WebSocket、TCP、KCP等）的连接实现都需要实现此接口
type Conn interface {
	// ID 返回连接的唯一标识符
	ID() int64
	// RemoteAddr 返回连接的远程地址
	RemoteAddr() net.Addr
	// LocalAddr 返回连接的本地地址
	LocalAddr() net.Addr
	// WriteTextMessage 向连接写入文本消息
	WriteTextMessage(msg []byte) error
	// WriteBinaryMessage 向连接写入二进制消息
	WriteBinaryMessage(msg []byte) error
	// Close 关闭连接
	Close()
}
