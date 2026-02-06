package pulse

import (
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

const (
	serverReadChunkSize = 32 * 1024
	serverReadPoolCap   = 64 * 1024
	serverReadPoolMax   = 256 * 1024
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, serverReadPoolCap)
		return &b
	},
}

type Conn struct {
	id   int64
	raw  net.Conn
	opts *options

	sendQ chan serverSendItem

	done chan struct{}

	closed atomic.Bool

	lastSeen atomic.Int64 // unix nano

	kv sync.Map
}

type serverSendItem struct {
	op  ws.OpCode
	msg []byte
}

func (s *Conn) ID() int64 { return s.id }

func (s *Conn) RemoteAddr() net.Addr { return s.raw.RemoteAddr() }

func (s *Conn) Set(key string, val any) { s.kv.Store(key, val) }

func (s *Conn) Get(key string) (any, bool) { return s.kv.Load(key) }

func (s *Conn) touch() { s.lastSeen.Store(time.Now().UnixNano()) }

func (s *Conn) LastSeen() time.Time { return time.Unix(0, s.lastSeen.Load()) }

func (s *Conn) Close() {
	if s.closed.CompareAndSwap(false, true) {
		s.trySendCloseFrame()
		close(s.done)
		_ = s.raw.Close()
	}
}

func (s *Conn) trySendCloseFrame() {
	payload := ws.NewCloseFrameBody(ws.StatusNormalClosure, "")
	if s.sendQ != nil {
		select {
		case s.sendQ <- serverSendItem{op: ws.OpClose, msg: payload}:
			return
		default:
		}
	}
	s.writeCloseFrameDirect(payload)
}

func (s *Conn) writeCloseFrameDirect(payload []byte) {
	if s.raw == nil {
		return
	}
	if s.opts != nil && s.opts.writeTimeout > 0 {
		_ = s.raw.SetWriteDeadline(time.Now().Add(s.opts.writeTimeout))
	}
	frame := ws.NewCloseFrame(payload)
	_ = ws.WriteFrame(s.raw, frame)
}

func (s *Conn) WriteBinary(msg []byte) error {
	return s.write(ws.OpBinary, msg)
}

func (s *Conn) WriteText(msg []byte) error {
	return s.write(ws.OpText, msg)
}

func (s *Conn) write(op ws.OpCode, msg []byte) error {
	if s.closed.Load() {
		return net.ErrClosed
	}

	// 上层可能复用 slice，必须拷贝，避免数据被改
	cp := make([]byte, len(msg))
	copy(cp, msg)

	switch s.opts.backpressure {
	case BackpressureBlock:
		select {
		case s.sendQ <- serverSendItem{op: op, msg: cp}:
			return nil
		case <-s.done:
			return net.ErrClosed
		}
	case BackpressureDrop:
		select {
		case s.sendQ <- serverSendItem{op: op, msg: cp}:
			return nil
		case <-s.done:
			return net.ErrClosed
		default:
			return nil
		}
	default: // Kick
		select {
		case s.sendQ <- serverSendItem{op: op, msg: cp}:
			return nil
		case <-s.done:
			return net.ErrClosed
		default:
			s.Close()
			return ErrBackpressure
		}
	}
}

func (s *Conn) writeLoop() {
	defer s.Close()

	w := wsutil.NewWriter(s.raw, ws.StateServerSide, ws.OpBinary) // w.ResetOp 会重置 op，所以这里先设置好

	for {
		select {
		case <-s.done:
			return
		case item := <-s.sendQ:
			if s.opts.writeTimeout > 0 {
				_ = s.raw.SetWriteDeadline(time.Now().Add(s.opts.writeTimeout))
			}
			switch item.op {
			case ws.OpBinary, ws.OpText, ws.OpClose:
				// ok
			default:
				continue
			}
			w.ResetOp(item.op)
			if _, err := w.Write(item.msg); err != nil {
				return
			}
			if err := w.Flush(); err != nil {
				return
			}
			if item.op == ws.OpClose {
				return
			}
		}
	}
}

func (s *Conn) readLoop() error {
	defer s.Close()

	controlHandler := wsutil.ControlFrameHandler(s.raw, ws.StateServerSide)
	rd := wsutil.Reader{
		Source:          s.raw,
		State:           ws.StateServerSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		MaxFrameSize:    s.opts.maxMessageSize,
		OnIntermediate:  controlHandler,
	}

	for {
		if s.opts.readTimeout > 0 {
			_ = s.raw.SetReadDeadline(time.Now().Add(s.opts.readTimeout))
		}

		hdr, err := rd.NextFrame()
		if err != nil {
			return err
		}

		// 控制帧交给ws处理,如: ping pong close
		if hdr.OpCode.IsControl() {
			if err = controlHandler(hdr, &rd); err != nil {
				return err
			}
			continue
		}

		if hdr.OpCode != ws.OpBinary && hdr.OpCode != ws.OpText {
			if err = rd.Discard(); err != nil {
				return err
			}
			continue
		}

		bp := bufPool.Get().(*[]byte)
		buf := (*bp)[:0]
		var tmp [serverReadChunkSize]byte

		for {
			n, err := rd.Read(tmp[:])
			if n > 0 {
				buf = append(buf, tmp[:n]...)
			}
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				bufPool.Put(bp)
				return err
			}
		}

		s.touch()

		if hdr.OpCode == ws.OpText && s.opts.onTextMessage != nil {
			s.opts.onTextMessage(s, buf)
		} else if hdr.OpCode == ws.OpBinary && s.opts.onBinaryMessage != nil {
			s.opts.onBinaryMessage(s, buf)
		}

		if cap(buf) > serverReadPoolMax {
			b := make([]byte, 0, serverReadPoolCap)
			*bp = b
		} else {
			*bp = buf[:0]
		}
		bufPool.Put(bp)
	}
}
