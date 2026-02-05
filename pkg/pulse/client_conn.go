package pulse

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func newClientConn(raw net.Conn, br *bufio.Reader, opts *clientOptions) *ClientConn {
	c := &ClientConn{
		raw:   raw,
		br:    br,
		opts:  opts,
		sendQ: make(chan clientSendItem, opts.sendQueueSize),
		done:  make(chan struct{}),
	}
	c.touch()
	c.touchPong()
	return c
}

const (
	clientReadChunkSize = 32 * 1024
	clientReadPoolCap   = 64 * 1024
	clientReadPoolMax   = 256 * 1024
)

var clientReadBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, clientReadPoolCap)
		return &b
	},
}

type ClientConn struct {
	raw  net.Conn
	br   *bufio.Reader
	opts *clientOptions

	sendQ chan clientSendItem

	done chan struct{}

	closed atomic.Bool

	lastSeen atomic.Int64 // unix nano
	lastPong atomic.Int64 // unix nano

	kv sync.Map
}

type clientSendItem struct {
	op  ws.OpCode
	msg []byte
}

func (c *ClientConn) RemoteAddr() net.Addr { return c.raw.RemoteAddr() }

func (c *ClientConn) Set(key string, val any) { c.kv.Store(key, val) }

func (c *ClientConn) Get(key string) (any, bool) { return c.kv.Load(key) }

func (c *ClientConn) touch() { c.lastSeen.Store(time.Now().UnixNano()) }

func (c *ClientConn) touchPong() { c.lastPong.Store(time.Now().UnixNano()) }

func (c *ClientConn) LastSeen() time.Time { return time.Unix(0, c.lastSeen.Load()) }

func (c *ClientConn) Close() {
	if c.closed.CompareAndSwap(false, true) {
		close(c.done)
		_ = c.raw.Close()
	}
}

func (c *ClientConn) WriteBinary(msg []byte) error {
	return c.write(ws.OpBinary, msg)
}

func (c *ClientConn) WriteText(msg []byte) error {
	return c.write(ws.OpText, msg)
}

func (c *ClientConn) write(op ws.OpCode, msg []byte) error {
	if c.closed.Load() {
		return net.ErrClosed
	}

	// 上层可能复用 slice，必须拷贝，避免数据被改
	cp := make([]byte, len(msg))
	copy(cp, msg)

	switch c.opts.backpressure {
	case BackpressureBlock:
		select {
		case c.sendQ <- clientSendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		}
	case BackpressureDrop:
		select {
		case c.sendQ <- clientSendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		default:
			return nil
		}
	default: // Kick
		select {
		case c.sendQ <- clientSendItem{op: op, msg: cp}:
			return nil
		case <-c.done:
			return net.ErrClosed
		default:
			c.Close()
			return ErrBackpressure
		}
	}
}

func (c *ClientConn) writeLoop() error {
	defer c.Close()

	w := wsutil.NewWriter(c.raw, ws.StateClientSide, ws.OpBinary)

	for {
		select {
		case <-c.done:
			return nil
		case item := <-c.sendQ:
			if c.opts.writeTimeout > 0 {
				_ = c.raw.SetWriteDeadline(time.Now().Add(c.opts.writeTimeout))
			}
			switch item.op {
			case ws.OpBinary, ws.OpText, ws.OpPing, ws.OpPong:
				// ok
			default:
				continue
			}
			w.ResetOp(item.op)
			if _, err := w.Write(item.msg); err != nil {
				return err
			}
			if err := w.Flush(); err != nil {
				return err
			}
		}
	}
}

func (c *ClientConn) readLoop() error {
	defer c.Close()
	defer func() {
		if c.br != nil {
			ws.PutReader(c.br)
		}
	}()

	controlHandler := wsutil.ControlFrameHandler(c.raw, ws.StateClientSide)

	source := io.Reader(c.raw)
	if c.br != nil {
		source = io.MultiReader(c.br, c.raw)
	}

	rd := wsutil.Reader{
		Source:          source,
		State:           ws.StateClientSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		OnIntermediate:  controlHandler,
	}

	for {
		if c.opts.readTimeout > 0 {
			_ = c.raw.SetReadDeadline(time.Now().Add(c.opts.readTimeout))
		}

		hdr, err := rd.NextFrame()
		if err != nil {
			return err
		}

		if hdr.OpCode.IsControl() {
			if hdr.OpCode == ws.OpPong {
				c.touchPong()
			}
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

		bp := clientReadBufPool.Get().(*[]byte)
		buf := (*bp)[:0]
		var tmp [clientReadChunkSize]byte

		for {
			n, err := rd.Read(tmp[:])
			if n > 0 {
				if c.opts.maxMessageSize > 0 && int64(len(buf)+n) > c.opts.maxMessageSize {
					clientReadBufPool.Put(bp)
					return wsutil.ErrFrameTooLarge
				}
				buf = append(buf, tmp[:n]...)
			}
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				clientReadBufPool.Put(bp)
				return err
			}
		}

		c.touch()

		if c.opts.onMessage != nil {
			c.opts.onMessage(c, hdr.OpCode, buf)
		}

		if cap(buf) > clientReadPoolMax {
			b := make([]byte, 0, clientReadPoolCap)
			*bp = b
		} else {
			*bp = buf[:0]
		}
		clientReadBufPool.Put(bp)
	}
}

func (c *ClientConn) pingLoop() error {
	if c.opts.pingInterval <= 0 {
		return nil
	}

	ticker := time.NewTicker(c.opts.pingInterval)
	defer ticker.Stop()

	timeout := c.opts.pingTimeout
	if timeout <= 0 {
		timeout = c.opts.pingInterval
	}

	for {
		select {
		case <-c.done:
			return nil
		case <-ticker.C:
		}

		sentAt := time.Now().UnixNano()
		if err := c.write(ws.OpPing, nil); err != nil {
			return err
		}

		if timeout <= 0 {
			continue
		}

		timer := time.NewTimer(timeout)
		select {
		case <-c.done:
			timer.Stop()
			return nil
		case <-timer.C:
		}

		if c.lastPong.Load() < sentAt {
			return ErrPingTimeout
		}
	}
}

func (c *ClientConn) run() error {
	defer c.Close()

	errCh := make(chan error, 3)
	go func() { errCh <- c.writeLoop() }()
	go func() { errCh <- c.readLoop() }()
	go func() { errCh <- c.pingLoop() }()

	return <-errCh
}
