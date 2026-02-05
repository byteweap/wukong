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

var serverReadBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, serverReadPoolCap)
		return &b
	},
}

type ServerConn struct {
	id   int64
	raw  net.Conn
	opts *serverOptions

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

func (s *ServerConn) ID() int64 { return s.id }

func (s *ServerConn) RemoteAddr() net.Addr { return s.raw.RemoteAddr() }

func (s *ServerConn) Set(key string, val any) { s.kv.Store(key, val) }

func (s *ServerConn) Get(key string) (any, bool) { return s.kv.Load(key) }

func (s *ServerConn) touch() { s.lastSeen.Store(time.Now().UnixNano()) }

func (s *ServerConn) LastSeen() time.Time { return time.Unix(0, s.lastSeen.Load()) }

func (s *ServerConn) Close() {
	if s.closed.CompareAndSwap(false, true) {
		close(s.done)
		_ = s.raw.Close()
	}
}

func (s *ServerConn) WriteBinary(msg []byte) error {
	return s.write(ws.OpBinary, msg)
}

func (s *ServerConn) WriteText(msg []byte) error {
	return s.write(ws.OpText, msg)
}

func (s *ServerConn) write(op ws.OpCode, msg []byte) error {
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

func (s *ServerConn) writeLoop() {
	defer s.Close()

	w := wsutil.NewWriter(s.raw, ws.StateServerSide, ws.OpBinary)

	for {
		select {
		case <-s.done:
			return
		case item := <-s.sendQ:
			if s.opts.writeTimeout > 0 {
				_ = s.raw.SetWriteDeadline(time.Now().Add(s.opts.writeTimeout))
			}
			if item.op != ws.OpBinary && item.op != ws.OpText {
				continue
			}
			w.ResetOp(item.op)
			if _, err := w.Write(item.msg); err != nil {
				return
			}
			if err := w.Flush(); err != nil {
				return
			}
		}
	}
}

func (s *ServerConn) readLoop() error {
	defer s.Close()

	controlHandler := wsutil.ControlFrameHandler(s.raw, ws.StateServerSide)
	rd := wsutil.Reader{
		Source:          s.raw,
		State:           ws.StateServerSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
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

		bp := serverReadBufPool.Get().(*[]byte)
		buf := (*bp)[:0]
		var tmp [serverReadChunkSize]byte

		for {
			n, err := rd.Read(tmp[:])
			if n > 0 {
				if s.opts.maxMessageSize > 0 && int64(len(buf)+n) > s.opts.maxMessageSize {
					serverReadBufPool.Put(bp)
					return wsutil.ErrFrameTooLarge
				}
				buf = append(buf, tmp[:n]...)
			}
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				serverReadBufPool.Put(bp)
				return err
			}
		}

		s.touch()

		if s.opts.onMessage != nil {
			s.opts.onMessage(s, hdr.OpCode, buf)
		}

		if cap(buf) > serverReadPoolMax {
			b := make([]byte, 0, serverReadPoolCap)
			*bp = b
		} else {
			*bp = buf[:0]
		}
		serverReadBufPool.Put(bp)
	}
}
