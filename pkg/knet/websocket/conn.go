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

	"github.com/byteweap/wukong/pkg/knet"
)

type Conn struct {
	opts         *Options
	id           int64
	raw          net.Conn
	writeQueue   chan writeMessage
	done         chan struct{}
	bufPool      sync.Pool
	lastPongTime int64 // 上次pong时间
}

type writeMessage struct {
	op   ws.OpCode
	data []byte
}

var _ knet.Conn = (*Conn)(nil)

func newConn(id int64, conn net.Conn, opts *Options) *Conn {

	return &Conn{
		opts:       opts,
		id:         id,
		raw:        conn,
		writeQueue: make(chan writeMessage, opts.WriteQueueSize),
		done:       make(chan struct{}),
		bufPool: sync.Pool{
			New: func() any {
				return make([]byte, 4*1024)
			},
		},
	}
}

func (c *Conn) ID() int64 {
	return c.id
}

func (c *Conn) LocalAddr() net.Addr {
	return c.raw.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.raw.RemoteAddr()
}

func (c *Conn) WriteBinaryMessage(msg []byte) error {
	select {
	case c.writeQueue <- writeMessage{ws.OpBinary, msg}:
		return nil
	default:
		return ErrWriteQueueFull
	}
}

func (c *Conn) WriteTextMessage(msg []byte) error {
	select {
	case c.writeQueue <- writeMessage{ws.OpText, msg}:
		return nil
	default:
		return ErrWriteQueueFull
	}
}

func (c *Conn) writeCloseFrame() {
	frame := ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, ""))
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write done-frame failed: %v", err))
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
		case <-c.done:
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

	defer c.close()

	for {
		select {
		case <-c.done:
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
	c.done <- struct{}{}
}

func (c *Conn) close() {

	close(c.done)
	close(c.writeQueue)

	if err := c.raw.Close(); err != nil {
		c.opts.handleError(fmt.Errorf("conn done failed: %v", err))
	}
}
