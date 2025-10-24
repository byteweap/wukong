package websocket

import (
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Conn struct {
	opts         *Options
	id           int64
	raw          net.Conn
	writeQueue   chan writeMessage
	close        chan struct{}
	closeOnce    sync.Once
	bufPool      sync.Pool
	lastPongTime int64 // 上次pong时间
}

type writeMessage struct {
	op   ws.OpCode
	data []byte
}

func newConn(id int64, conn net.Conn, opts *Options) *Conn {

	c := &Conn{
		opts:       opts,
		id:         id,
		raw:        conn,
		writeQueue: make(chan writeMessage, opts.WriteQueueSize),
		close:      make(chan struct{}),
		bufPool: sync.Pool{
			New: func() any {
				return make([]byte, 4*1024)
			},
		},
	}
	return c
}

func (c *Conn) WriteMessage(op ws.OpCode, msg []byte) error {
	select {
	case c.writeQueue <- writeMessage{op, msg}:
		return nil
	default:
		return ErrWriteQueueFull
	}
}

func (c *Conn) writeCloseFrame() {
	frame := ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, ""))
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write close-frame failed: %v", err))
	}
}

func (c *Conn) writePongFrame() {
	frame := ws.NewPongFrame(nil)
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write pong-frame failed: %v", err))
	}
}

func (c *Conn) writePingFrame() {
	frame := ws.NewPingFrame(nil)
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write ping-frame failed: %v", err))
	}
}

func (c *Conn) writePump() {

	ticker := time.NewTicker(c.opts.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.close:
			return
		case <-ticker.C:
			c.writePingFrame() // heartbeat
		case msg, ok := <-c.writeQueue:
			if !ok {
				return
			}
			// 高性能 writer（复用）
			w := wsutil.GetWriter(c.raw, ws.StateServerSide, msg.op, len(msg.data))
			w.Reset(c.raw, ws.StateServerSide, msg.op)

			// 可选：设置写超时，避免网络卡死
			if err := c.raw.SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout)); err != nil {
				c.opts.handleError(fmt.Errorf("set write deadline failed: %v", err))
			}

			if _, err := w.Write(msg.data); err != nil {
				c.Close()
				wsutil.PutWriter(w)
				c.opts.handleError(fmt.Errorf("write message failed: %v", err))
				return
			}

			if err := w.Flush(); err != nil {
				c.Close()
				wsutil.PutWriter(w)
				c.opts.handleError(fmt.Errorf("write flush failed: %v", err))
				return
			}
			wsutil.PutWriter(w)
		}
	}
}

func (c *Conn) readPump() {

	defer c.Close()

	for {
		select {
		case <-c.close:
			return
		default:
			header, err := ws.ReadHeader(c.raw)
			if err != nil {
				return
			}
			switch header.OpCode {
			case ws.OpContinuation:
				continue
			case ws.OpPing:
				c.writePongFrame()
			case ws.OpPong:
				atomic.StoreInt64(&c.lastPongTime, time.Now().UnixNano())
			case ws.OpClose:
				c.writeCloseFrame()
				return
			case ws.OpText, ws.OpBinary:
				if c.opts.MaxMessageSize > 0 && header.Length > c.opts.MaxMessageSize {
					c.opts.handleError(fmt.Errorf("message too large: %d > %d", header.Length, c.opts.MaxMessageSize))
					return
				}
				payload := c.bufPool.Get().([]byte)[:header.Length]
				if _, err = io.ReadFull(c.raw, payload); err != nil {
					c.bufPool.Put(payload)
					c.opts.handleError(fmt.Errorf("read message read full failed: %v", err))
					return
				}
				if header.OpCode == ws.OpText {
					c.opts.handleMessage(c, payload)
				} else if header.OpCode == ws.OpBinary {
					c.opts.handleBinaryMessage(c, payload)
				}
				c.bufPool.Put(payload)
			}
		}
	}
}

func (c *Conn) Close() {
	var err error
	c.closeOnce.Do(func() {
		close(c.close)
		err = c.raw.Close()
	})
	if err != nil {
		c.opts.handleError(fmt.Errorf("conn close failed: %v", err))
	}
}
