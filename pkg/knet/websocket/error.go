package websocket

import "errors"

var (
	ErrInvalidOpCode  = errors.New("invalid opcode")
	ErrInvalidState   = errors.New("invalid state")
	ErrConnClosed     = errors.New("connection closed")
	ErrWriteQueueFull = errors.New("write queue full")
	ErrMaxConns       = errors.New("max connections reached")
)
