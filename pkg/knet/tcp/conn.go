package tcp

import (
	"net"

	"github.com/byteweap/wukong/pkg/knet"
)

type tcpConn struct {
	id  int64
	raw *net.TCPConn
}

var _ knet.Conn = (*tcpConn)(nil)

func newConn(id int64, raw *net.TCPConn) *tcpConn {
	return &tcpConn{
		id:  id,
		raw: raw,
	}
}

func (t *tcpConn) ID() int64 {
	//TODO implement me
	panic("implement me")
}

func (t *tcpConn) RemoteAddr() net.Addr {
	return t.raw.RemoteAddr()
}

func (t *tcpConn) LocalAddr() net.Addr {
	return t.raw.LocalAddr()
}

func (t *tcpConn) WriteTextMessage(msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *tcpConn) WriteBinaryMessage(msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *tcpConn) Close() {
	t.raw.Close()
}
