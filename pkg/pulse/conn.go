package pulse

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Conn struct {
	id   int64
	uid  atomic.Int64
	raw  net.Conn
	opts *options

	sendQ chan []byte

	ctx    context.Context
	cancel context.CancelFunc

	closed atomic.Bool

	lastSeen atomic.Int64 // unix nano

	kv sync.Map
}

func (s *Conn) ID() int64 { return s.id }

func (s *Conn) UID() int64 { return s.uid.Load() }

func (s *Conn) BindUID(uid int64) {
	s.uid.Store(uid)
}

func (s *Conn) RemoteAddr() net.Addr { return s.raw.RemoteAddr() }

func (s *Conn) Set(key string, val any) { s.kv.Store(key, val) }

func (s *Conn) Get(key string) (any, bool) { return s.kv.Load(key) }

func (s *Conn) touch() { s.lastSeen.Store(time.Now().UnixNano()) }

func (s *Conn) LastSeen() time.Time { return time.Unix(0, s.lastSeen.Load()) }

func (s *Conn) Close() {
	if s.closed.CompareAndSwap(false, true) {
		s.cancel()
		_ = s.raw.Close()
	}
}

func (s *Conn) Write(msg []byte) error {
	if s.closed.Load() {
		return net.ErrClosed
	}

	// 上层可能复用 slice，必须拷贝，避免数据被改
	cp := make([]byte, len(msg))
	copy(cp, msg)

	switch s.opts.Backpressure {
	case BackpressureBlock:
		select {
		case s.sendQ <- cp:
			return nil
		case <-s.ctx.Done():
			return context.Canceled
		}
	case BackpressureDrop:
		select {
		case s.sendQ <- cp:
			return nil
		default:
			return nil
		}
	default: // Kick
		select {
		case s.sendQ <- cp:
			return nil
		default:
			s.Close()
			return context.DeadlineExceeded
		}
	}
}

func (s *Conn) writeLoop() {
	defer s.Close()

	for {
		select {
		case <-s.ctx.Done():
			return
		case msg := <-s.sendQ:
			if s.opts.WriteTimeout > 0 {
				_ = s.raw.SetWriteDeadline(time.Now().Add(s.opts.WriteTimeout))
			}
			// 二进制帧（游戏更常用）
			if err := wsutil.WriteServerBinary(s.raw, msg); err != nil {
				return
			}
		}
	}
}

func (s *Conn) readLoop() error {
	defer s.Close()

	for {
		if s.opts.ReadTimeout > 0 {
			_ = s.raw.SetReadDeadline(time.Now().Add(s.opts.ReadTimeout))
		}
		// 读取完整 message(简化：data + opcode)
		data, op, err := wsutil.ReadClientData(s.raw)
		if err != nil {
			return err
		}
		s.touch()

		// 限制消息大小（防御）
		if s.opts.MaxMessageSize > 0 && len(data) > s.opts.MaxMessageSize {
			return wsutil.ErrFrameTooLarge
		}

		switch op {
		case ws.OpBinary, ws.OpText:
			s.opts.onMessage(s, data)
		case ws.OpPing:
			_ = ws.WriteFrame(s.raw, ws.NewPongFrame(nil))
		case ws.OpClose:
			return nil
		default:
			// ignore
		}
	}
}
