package pulse

import (
	"errors"
	"fmt"
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

	sendQ chan sendItem

	done chan struct{}

	closed atomic.Bool

	lastSeen atomic.Int64 // unix nano

	kv sync.Map
}

type sendItem struct {
	op  ws.OpCode
	msg []byte
}

func (c *Conn) ID() int64 { return c.id }

func (c *Conn) RemoteAddr() net.Addr { return c.raw.RemoteAddr() }

func (c *Conn) Set(key string, val any) { c.kv.Store(key, val) }

func (c *Conn) Get(key string) (any, bool) { return c.kv.Load(key) }

func (c *Conn) touch() { c.lastSeen.Store(time.Now().UnixNano()) }

func (c *Conn) LastSeen() time.Time { return time.Unix(0, c.lastSeen.Load()) }

func (c *Conn) Close() {
	if c.closed.CompareAndSwap(false, true) {
		c.trySendCloseFrame()
		close(c.done)
		_ = c.raw.Close()
	}
}

func (c *Conn) trySendCloseFrame() {
	payload := ws.NewCloseFrameBody(ws.StatusNormalClosure, "")
	if c.sendQ != nil {
		select {
		case c.sendQ <- sendItem{op: ws.OpClose, msg: payload}:
			return
		default:
		}
	}
	c.writeCloseFrameDirect(payload)
}

func (c *Conn) writeCloseFrameDirect(payload []byte) {
	if c.raw == nil {
		return
	}
	if c.opts != nil && c.opts.writeTimeout > 0 {
		_ = c.raw.SetWriteDeadline(time.Now().Add(c.opts.writeTimeout))
	}
	frame := ws.NewCloseFrame(payload)
	_ = ws.WriteFrame(c.raw, frame)
}

func (c *Conn) WriteBinary(msg []byte) error {
	return c.write(ws.OpBinary, msg)
}

func (c *Conn) WriteText(msg []byte) error {
	return c.write(ws.OpText, msg)
}

func (c *Conn) write(op ws.OpCode, msg []byte) error {
	if c.closed.Load() {
		return net.ErrClosed
	}

	// 上层可能复用 slice，必须拷贝，避免数据被改
	cp := make([]byte, len(msg))
	copy(cp, msg)

	switch c.opts.backpressure {
	case BackpressureBlock:
		select {
		case c.sendQ <- sendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		}
	case BackpressureDrop:
		select {
		case c.sendQ <- sendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		default:
			return nil
		}
	default: // Kick
		select {
		case c.sendQ <- sendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		default:
			c.Close()
			return ErrBackpressure
		}
	}
}

func (c *Conn) writeLoop() {
	defer c.Close()

	tick := time.NewTicker(c.opts.pingInterval)
	defer tick.Stop()

	var err error
	defer func() {
		if err != nil && c.opts.errorHandler != nil {
			c.opts.errorHandler(c, fmt.Errorf("write loop error: %w", err))
		}
	}()

	w := wsutil.NewWriter(c.raw, ws.StateServerSide, 0)

	for {
		select {
		case <-c.done:
			return
		case <-tick.C:

			w.ResetOp(ws.OpPing)
			if _, err = w.Write(nil); err != nil {
				return
			}
			if err = w.Flush(); err != nil {
				return
			}

		case item := <-c.sendQ:

			if c.opts.writeTimeout > 0 {
				if err = c.raw.SetWriteDeadline(time.Now().Add(c.opts.writeTimeout)); err != nil {
					return
				}
			}
			switch item.op {
			case ws.OpBinary, ws.OpText, ws.OpClose:
				// ok
			default:
				continue
			}
			w.ResetOp(item.op)
			if _, err = w.Write(item.msg); err != nil {
				return
			}
			if err = w.Flush(); err != nil {
				return
			}
			if item.op == ws.OpClose {
				return
			}
		}
	}

}

func (c *Conn) readLoop() {
	defer c.Close()

	var err error
	defer func() {
		if err != nil && c.opts.errorHandler != nil {
			c.opts.errorHandler(c, fmt.Errorf("read loop error: %w", err))
		}
	}()

	ctl := wsutil.ControlFrameHandler(c.raw, ws.StateServerSide)
	rd := wsutil.Reader{
		Source:          c.raw,
		State:           ws.StateServerSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		MaxFrameSize:    c.opts.maxMessageSize,
		OnIntermediate:  ctl,
	}

	for {
		if c.opts.readTimeout > 0 {
			if err = c.raw.SetReadDeadline(time.Now().Add(c.opts.readTimeout)); err != nil {
				return
			}
		}

		hdr, err := rd.NextFrame()
		if err != nil {
			return
		}

		// 控制帧交给ws处理,如: ping pong close
		if hdr.OpCode.IsControl() {
			if err = ctl(hdr, &rd); err != nil {
				return
			}
			continue
		}

		if hdr.OpCode != ws.OpBinary && hdr.OpCode != ws.OpText {
			if err = rd.Discard(); err != nil {
				return
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
				return
			}
		}

		c.touch()

		if hdr.OpCode == ws.OpText && c.opts.onTextMessage != nil {
			c.opts.onTextMessage(c, buf)
		} else if hdr.OpCode == ws.OpBinary && c.opts.onBinaryMessage != nil {
			c.opts.onBinaryMessage(c, buf)
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
