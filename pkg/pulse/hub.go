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

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c, _ := range h.cs {
		for _, filter := range filters {
			if !filter(c) {
				return
			}
		}
		_ = c.WriteBinary(msg)
	}
}

func (h *hub) broadcastText(msg []byte, filters ...func(conn *Conn) bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c, _ := range h.cs {
		for _, filter := range filters {
			if !filter(c) {
				return
			}
		}
		_ = c.WriteText(msg)
	}
}

func (h *hub) allocate(parent context.Context, opts *options, conn net.Conn) {

	ctx, cancel := context.WithCancel(parent)

	id := h.nextID.Add(1)
	s := &Conn{
		id:     id,
		raw:    conn,
		sendQ:  make(chan sendItem, opts.sendQueueSize),
		ctx:    ctx,
		cancel: cancel,
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
}
