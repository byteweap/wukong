package websocket

type (
	StartHandler      func(addr, pattern string)
	StopHandler       func()
	ConnectHandler    func(conn *Conn)
	DisconnectHandler func(conn *Conn)
	MessageHandler    func(conn *Conn, msg []byte)
	ErrorHandler      func(err error)
)
