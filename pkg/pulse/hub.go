package pulse

import (
	"net"
	"sync"
	"sync/atomic"
)

type hub struct {
	mu     sync.RWMutex
	cs     map[*Conn]struct{}
	nextID atomic.Int64
	open   atomic.Bool
}

func newHub() *hub {
	h := &hub{
		cs:     make(map[*Conn]struct{}),
		nextID: atomic.Int64{},
	}
	h.nextID.Store(0)
	h.open.Store(true)
	return h
}

func (h *hub) closed() bool {
	return !h.open.Load()
}

func (h *hub) register(s *Conn) {
	h.mu.Lock()
	h.cs[s] = struct{}{}
	h.mu.Unlock()
}

func (h *hub) unregister(s *Conn) {
	h.mu.Lock()
	delete(h.cs, s)
	h.mu.Unlock()
}

func (h *hub) broadcastBinary(msg []byte, filters ...func(conn *Conn) bool) {

	for _, c := range h.all() {
		for _, filter := range filters {
			if !filter(c) {
				break
			}
		}
		_ = c.WriteBinary(msg)
	}
}

func (h *hub) broadcastText(msg []byte, filters ...func(conn *Conn) bool) {

	for _, c := range h.all() {
		for _, filter := range filters {
			if !filter(c) {
				break
			}
		}
		_ = c.WriteText(msg)
	}
}

func (h *hub) allocate(opts *options, conn net.Conn) error {

	if h.closed() {
		return ErrClosed
	}

	id := h.nextID.Add(1)
	s := &Conn{
		id:    id,
		opts:  opts,
		raw:   conn,
		sendQ: make(chan sendItem, opts.sendQueueSize),
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
	s.readLoop()

	// 清理
	h.unregister(s)

	if opts.onDisconnect != nil {
		opts.onDisconnect(s)
	}

	s.Close()
	return nil
}

func (h *hub) all() []*Conn {
	h.mu.RLock()
	cs := make([]*Conn, 0, len(h.cs))
	for c := range h.cs {
		cs = append(cs, c)
	}
	h.mu.RUnlock()
	return cs
}

func (h *hub) close() {
	if !h.open.CompareAndSwap(true, false) {
		return
	}
	for _, c := range h.all() {
		c.Close()
	}
}
