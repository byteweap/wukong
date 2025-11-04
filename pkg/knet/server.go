package knet

type (
	StartHandler       func(addr, pattern string)
	StopHandler        func(err error)
	ConnectHandler     func(conn Conn)
	ConnMessageHandler func(conn Conn, msg []byte)
	ErrorHandler       func(err error)
)

type Server interface {
	// Addr return the server address.
	Addr() string
	// Protocol return protocol name.
	Protocol() string
	// Start starts the server.
	Start()
	// Shutdown graceful stop the server.
	Shutdown()
	// OnStart set the handler to be called when the server starts.
	OnStart(StartHandler)
	// OnStop set the handler to be called when the server stops.
	OnStop(StopHandler)
	// OnConnect set the handler to be called when a new connection is established.
	OnConnect(ConnectHandler)
	// OnDisconnect set the handler to be called when a connection is closed.
	OnDisconnect(ConnectHandler)
	// OnTextMessage set the handler to be called when a text message is received.
	OnTextMessage(ConnMessageHandler)
	// OnBinaryMessage set the handler to be called when a binary message is received.
	OnBinaryMessage(ConnMessageHandler)
	// OnError set the handler to be called when find error.
	OnError(ErrorHandler)
}
