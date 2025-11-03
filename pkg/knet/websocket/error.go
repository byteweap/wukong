package websocket

import "errors"

var (
	ErrInvalidOpCode  = errors.New("invalid opcode")
	ErrInvalidState   = errors.New("invalid state")
	ErrHubClosed      = errors.New("hub closed")
	ErrWriteQueueFull = errors.New("write queue full")
	ErrMaxConns       = errors.New("max connections reached")
)
