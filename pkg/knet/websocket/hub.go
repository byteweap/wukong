package websocket

import (
	"net"
	"sync"
	"sync/atomic"
)

type hub struct {
	opts       *Options
	nextConnID int64
	open       atomic.Bool // 是否打开

	mux        sync.RWMutex
	totalConns atomic.Int64       // 总链接数
	conns      map[*Conn]struct{} // 所有链接
}

func newHub(opts *Options) *hub {
	h := &hub{
		opts:       opts,
		nextConnID: 0,
		conns:      make(map[*Conn]struct{}),
	}
	h.open.Store(true)
	h.totalConns.Store(0)
	return h
}

func (h *hub) closed() bool {
	return !h.open.Load()
}

func (h *hub) shutdown() error {

	if h.closed() {
		return ErrHubClosed
	}
	h.open.Store(false)

	h.mux.RLock()
	for conn, _ := range h.conns {
		conn.Close()
	}
	h.mux.RUnlock()

	clear(h.conns)

	return nil
}

func (h *hub) allocate(netConn net.Conn) error {
	if h.closed() {
		return ErrHubClosed
	}
	if h.totalConns.Load() >= int64(h.opts.MaxConnections) {
		return ErrMaxConns
	}

	nextConnID := atomic.AddInt64(&h.nextConnID, 1)
	conn := newConn(nextConnID, netConn, h.opts)

	h.mux.Lock()
	h.conns[conn] = struct{}{}
	h.mux.Unlock()

	h.totalConns.Add(1)

	go conn.writePump()
	conn.readPump()

	h.totalConns.Add(-1)

	return nil
}
