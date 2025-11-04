package websocket

import (
	"math"
	"net"
	"sync"
	"sync/atomic"
)

type hub struct {
	opts       *Options
	nextConnID int64
	open       atomic.Bool // 是否打开

	mux     sync.RWMutex
	connNum atomic.Int64       // 总链接数
	conns   map[*Conn]struct{} // 所有链接
}

func newHub(opts *Options) *hub {
	h := &hub{
		opts:       opts,
		nextConnID: 0,
		conns:      make(map[*Conn]struct{}),
	}
	h.open.Store(true)
	h.connNum.Store(0)
	return h
}

func (h *hub) closed() bool {
	return !h.open.Load()
}

// graceful stop hub
func (h *hub) shutdown() error {

	if h.closed() {
		return ErrHubClosed
	}
	h.open.Store(false)

	h.mux.RLock()
	for conn := range h.conns {
		conn.Close()
	}
	h.mux.RUnlock()

	clear(h.conns)

	return nil
}

// next connection id
func (h *hub) nextId() int64 {

	atomic.AddInt64(&h.nextConnID, 1)

	if atomic.LoadInt64(&h.nextConnID) == math.MaxInt64 {
		atomic.StoreInt64(&h.nextConnID, 1)
	}
	return atomic.LoadInt64(&h.nextConnID)
}

// allocate connection
func (h *hub) allocate(netConn net.Conn) error {
	if h.closed() {
		return ErrHubClosed
	}
	if h.connNum.Load() >= int64(h.opts.MaxConnections) {
		return ErrMaxConns
	}
	id := h.nextId()
	conn := newConn(id, netConn, h.opts)

	h.mux.Lock()
	h.conns[conn] = struct{}{}
	h.mux.Unlock()
	h.connNum.Add(1)

	h.opts.handleConnect(conn) // connect

	go conn.writePump()
	conn.readPump()

	h.mux.Lock()
	delete(h.conns, conn)
	h.mux.Unlock()
	h.connNum.Add(-1)

	h.opts.handleDisconnect(conn) // disconnect

	return nil
}
