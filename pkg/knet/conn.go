package knet

import "net"

type Conn interface {
	// ID returns the connection ID.
	ID() int64
	// RemoteAddr returns the remote address of the connection.
	RemoteAddr() net.Addr
	// LocalAddr returns the local address of the connection.
	LocalAddr() net.Addr
	// WriteTextMessage writes a text message to the connection.
	WriteTextMessage(msg []byte) error
	// WriteBinaryMessage writes a binary message to the connection.
	WriteBinaryMessage(msg []byte) error
	// Close closes the connection.
	Close()
}
