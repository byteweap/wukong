package pulse

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
)

type hub struct {
	mu     sync.RWMutex
	cs     map[*Conn]struct{}
	nextID atomic.Int64
}

func newHub() *hub {
	h := &hub{
		cs:     make(map[*Conn]struct{}),
		nextID: atomic.Int64{},
	}
	h.nextID.Store(0)
	return h
}

func (m *hub) add(s *Conn) {
	m.mu.Lock()
	m.cs[s] = struct{}{}
	m.mu.Unlock()
}

func (m *hub) remove(s *Conn) {
	m.mu.Lock()
	delete(m.cs, s)
	m.mu.Unlock()
}

func (m *hub) rangeAll(fn func(*Conn)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for c, _ := range m.cs {
		fn(c)
	}
}

func (m *hub) broadcast(msg []byte) {
	m.rangeAll(func(s *Conn) {
		_ = s.Write(msg)
	})
}

func (h *hub) allocate(parent context.Context, opts *options, conn net.Conn) {

	ctx, cancel := context.WithCancel(parent)

	id := h.nextID.Add(1)
	s := &Conn{
		id:     id,
		raw:    conn,
		sendQ:  make(chan []byte, opts.SendQueueSize),
		ctx:    ctx,
		cancel: cancel,
	}
	s.touch()

	h.add(s)

	if opts.onConnect != nil {
		opts.onConnect(s)
	}

	// 写协程
	go s.writeLoop()

	// 读循环（当前协程）
	readErr := s.readLoop()

	// 清理
	h.remove(s)
	s.Close()

	if readErr != nil && opts.onError != nil {
		opts.onError(s, readErr)
	}
	if opts.onDisconnect != nil {
		opts.onDisconnect(s, readErr)
	}
}
