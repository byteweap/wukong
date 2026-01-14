package network

// 服务器事件处理器类型定义
type (

	// StartHandler 服务器启动时的回调函数
	StartHandler func(addr, pattern string)

	// StopHandler 服务器停止时的回调函数
	StopHandler func()

	// ConnectHandler 连接建立或断开时的回调函数
	ConnectHandler func(conn Conn)

	// ConnMessageHandler 接收到消息时的回调函数
	ConnMessageHandler func(conn Conn, msg []byte)

	// ErrorHandler 错误发生时的回调函数
	ErrorHandler func(err error)
)

// Server 定义网络服务器的核心接口
// 所有网络协议（WebSocket、TCP、KCP等）的服务器实现都需要实现此接口
type Server interface {

	// Addr 返回服务器监听地址
	Addr() string

	// Protocol 返回协议名称
	Protocol() string

	// Start 启动服务器
	Start()

	// OnStart 设置服务器启动时的回调函数
	OnStart(StartHandler)

	// OnStop 设置服务器停止时的回调函数
	OnStop(StopHandler)

	// OnConnect 设置新连接建立时的回调函数
	OnConnect(ConnectHandler)

	// OnDisconnect 设置连接关闭时的回调函数
	OnDisconnect(ConnectHandler)

	// OnTextMessage 设置接收到文本消息时的回调函数
	OnTextMessage(ConnMessageHandler)

	// OnBinaryMessage 设置接收到二进制消息时的回调函数
	OnBinaryMessage(ConnMessageHandler)

	// OnError 设置错误发生时的回调函数
	OnError(ErrorHandler)

	// Stop 优雅关闭服务器
	Stop()
}
