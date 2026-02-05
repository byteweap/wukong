package pulse

import (
	"net"
	"sync"
	"sync/atomic"
)

type serverHub struct {
	mu     sync.RWMutex
	cs     map[*ServerConn]struct{}
	nextID atomic.Int64
	open   atomic.Bool
}

func newServerHub() *serverHub {
	h := &serverHub{
		cs:     make(map[*ServerConn]struct{}),
		nextID: atomic.Int64{},
	}
	h.nextID.Store(0)
	h.open.Store(true)
	return h
}

func (h *serverHub) closed() bool {
	return !h.open.Load()
}

func (h *serverHub) register(s *ServerConn) {
	h.mu.Lock()
	h.cs[s] = struct{}{}
	h.mu.Unlock()
}

func (h *serverHub) unregister(s *ServerConn) {
	h.mu.Lock()
	delete(h.cs, s)
	h.mu.Unlock()
}

func (h *serverHub) broadcastBinary(msg []byte, filters ...func(conn *ServerConn) bool) {

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.cs {
		for _, filter := range filters {
			if !filter(c) {
				break
			}
		}
		_ = c.WriteBinary(msg)
	}
}

func (h *serverHub) broadcastText(msg []byte, filters ...func(conn *ServerConn) bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.cs {
		for _, filter := range filters {
			if !filter(c) {
				break
			}
		}
		_ = c.WriteText(msg)
	}
}

func (h *serverHub) allocate(opts *serverOptions, conn net.Conn) error {

	if h.closed() {
		return ErrClosed
	}

	id := h.nextID.Add(1)
	s := &ServerConn{
		id:    id,
		opts:  opts,
		raw:   conn,
		sendQ: make(chan serverSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}
	s.touch()

	h.register(s)

	if opts.onConnect != nil {
		opts.onConnect(s)
	}

	// 写协程
	go s.writeLoop()

	// 读循环（当前协程）
	readErr := s.readLoop()

	// 清理
	h.unregister(s)
	s.Close()

	if readErr != nil && opts.onError != nil {
		opts.onError(s, readErr)
	}
	if opts.onDisconnect != nil {
		opts.onDisconnect(s, readErr)
	}

	return nil
}

func (h *serverHub) closeAll() {
	h.mu.RLock()
	conns := make([]*ServerConn, 0, len(h.cs))
	for c := range h.cs {
		conns = append(conns, c)
	}
	h.mu.RUnlock()

	for _, c := range conns {
		c.Close()
	}
}

func (h *serverHub) close() {
	if !h.open.CompareAndSwap(true, false) {
		return
	}
	h.closeAll()
}
