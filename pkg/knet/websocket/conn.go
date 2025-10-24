package websocket

import (
	"io"
	"log"
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
	_ = ws.WriteFrame(c.raw, frame)
}

func (c *Conn) writePongFrame() {
	frame := ws.NewPongFrame(nil)
	_ = ws.WriteFrame(c.raw, frame)
}

func (c *Conn) writePingFrame() {
	frame := ws.NewPingFrame(nil)
	_ = ws.WriteFrame(c.raw, frame)
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

			_ = c.raw.SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout)) // 可选：设置写超时，避免网络卡死

			if _, err := w.Write(msg.data); err != nil {
				_ = c.Close()
				wsutil.PutWriter(w)
				return
			}
			if err := w.Flush(); err != nil {
				_ = c.Close()
				wsutil.PutWriter(w)
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
			case ws.OpClose:
				c.writeCloseFrame()
				return
			case ws.OpPing:
				c.writePongFrame()
			case ws.OpPong:
				atomic.StoreInt64(&c.lastPongTime, time.Now().UnixNano())
			case ws.OpText, ws.OpBinary:

				if c.opts.MaxMessageSize > 0 && header.Length > int64(c.opts.MaxMessageSize) {
					log.Printf("message too large: %d > %d", header.Length, c.opts.MaxMessageSize)
					return
				}
				payload := c.bufPool.Get().([]byte)[:header.Length]
				if _, err = io.ReadFull(c.raw, payload); err != nil {
					c.bufPool.Put(payload)
					return
				}
				if c.opts.onMessageHandler != nil {
					c.opts.onMessageHandler(c, header.OpCode, payload)
				}
				c.bufPool.Put(payload)

			default:
				continue
			}
		}
	}
}

func (c *Conn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		close(c.close)
		err = c.raw.Close()
	})
	return err
}
