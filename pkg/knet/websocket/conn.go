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

// Conn 实现WebSocket连接的具体结构体
// 实现knet.Conn接口，管理单个WebSocket连接的生命周期和消息处理
type Conn struct {
	opts         *Options          // 连接配置选项
	id           int64             // 连接唯一标识符
	raw          net.Conn          // 底层网络连接
	writeQueue   chan writeMessage // 写入消息队列，实现异步写入
	done         chan struct{}     // 用于通知关闭的通道
	bufPool      sync.Pool         // 缓冲区池，减少内存分配
	lastPongTime int64             // 上次收到pong消息的时间（用于心跳检测）
}

// writeMessage 定义写入队列中的消息结构
type writeMessage struct {
	op   ws.OpCode // WebSocket操作码
	data []byte    // 消息数据
}

// 确保Conn实现了knet.Conn接口
var _ knet.Conn = (*Conn)(nil)

// newConn 创建新的WebSocket连接实例
func newConn(id int64, conn net.Conn, opts *Options) *Conn {
	return &Conn{
		opts:       opts,
		id:         id,
		raw:        conn,
		writeQueue: make(chan writeMessage, opts.WriteQueueSize), // 创建带缓冲的写入队列
		done:       make(chan struct{}),
		bufPool: sync.Pool{ // 初始化缓冲区池，默认4KB缓冲区
			New: func() any {
				return make([]byte, 4*1024)
			},
		},
	}
}

// ID 返回连接的唯一标识符
func (c *Conn) ID() int64 {
	return c.id
}

// LocalAddr 返回连接的本地地址
func (c *Conn) LocalAddr() net.Addr {
	return c.raw.LocalAddr()
}

// RemoteAddr 返回连接的远程地址
func (c *Conn) RemoteAddr() net.Addr {
	return c.raw.RemoteAddr()
}

// WriteBinaryMessage 向连接写入二进制消息
func (c *Conn) WriteBinaryMessage(msg []byte) error {
	select {
	case c.writeQueue <- writeMessage{ws.OpBinary, msg}:
		return nil
	default:
		return ErrWriteQueueFull // 队列满，返回错误
	}
}

// WriteTextMessage 向连接写入文本消息
func (c *Conn) WriteTextMessage(msg []byte) error {
	select {
	case c.writeQueue <- writeMessage{ws.OpText, msg}:
		return nil
	default:
		return ErrWriteQueueFull // 队列满，返回错误
	}
}

// writeCloseFrame 发送关闭帧
func (c *Conn) writeCloseFrame() {
	frame := ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, ""))
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write done-frame failed: %v", err))
	}
}

// writePongFrame 发送pong帧
func (c *Conn) writePongFrame() {
	frame := ws.NewPongFrame(nil)
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write pong-frame failed: %v", err))
	}
}

// writePingFrame 发送ping帧
func (c *Conn) writePingFrame() {
	frame := ws.NewPingFrame(nil)
	if err := ws.WriteFrame(c.raw, frame); err != nil {
		c.opts.handleError(fmt.Errorf("write ping-frame failed: %v", err))
	}
}

// writePump 处理连接的异步写入操作
func (c *Conn) writePump() {
	// 创建心跳定时器
	ticker := time.NewTicker(c.opts.PingInterval)
	defer ticker.Stop() // 确保定时器被停止

	// 主循环处理写入操作
	for {
		select {
		case <-c.done: // 连接关闭信号
			return
		case <-ticker.C: // 心跳触发
			c.writePingFrame() // 发送ping帧保持连接活跃
		case msg, ok := <-c.writeQueue: // 从写入队列获取消息
			if !ok { // 队列已关闭
				return
			}
			// 从对象池获取高性能writer（减少内存分配）
			w := wsutil.GetWriter(c.raw, ws.StateServerSide, msg.op, len(msg.data))
			w.Reset(c.raw, ws.StateServerSide, msg.op)

			// 设置写超时，避免网络卡死导致goroutine泄漏
			if err := c.raw.SetWriteDeadline(time.Now().Add(c.opts.WriteTimeout)); err != nil {
				c.opts.handleError(fmt.Errorf("set write deadline failed: %v", err))
			}

			// 写入消息数据
			if _, err := w.Write(msg.data); err != nil {
				c.Close()           // 写入失败，关闭连接
				wsutil.PutWriter(w) // 归还writer到对象池
				c.opts.handleError(fmt.Errorf("write message failed: %v", err))
				return
			}

			// 刷新缓冲区确保数据发送
			if err := w.Flush(); err != nil {
				c.Close()           // 刷新失败，关闭连接
				wsutil.PutWriter(w) // 归还writer到对象池
				c.opts.handleError(fmt.Errorf("write flush failed: %v", err))
				return
			}
			wsutil.PutWriter(w) // 归还writer到对象池
		}
	}
}

// readPump 处理连接的读取操作
func (c *Conn) readPump() {
	defer c.close() // 确保连接正确关闭

	// 主循环处理读取操作
	for {
		select {
		case <-c.done: // 连接关闭信号
			return
		default:
			// 读取WebSocket帧头
			header, err := ws.ReadHeader(c.raw)
			if err != nil { // 读取失败，可能是连接断开
				return
			}

			// 根据操作码处理不同类型的帧
			switch header.OpCode {
			case ws.OpContinuation: // 继续帧，暂时忽略
				continue
			case ws.OpPing: // ping帧，回复pong帧
				c.writePongFrame()
			case ws.OpPong: // pong帧，更新最后收到pong的时间
				atomic.StoreInt64(&c.lastPongTime, time.Now().UnixNano())
			case ws.OpClose: // 关闭帧，回复关闭帧并结束连接
				c.writeCloseFrame()
				return
			case ws.OpText, ws.OpBinary: // 文本或二进制消息帧
				// 检查消息大小是否超过限制
				if c.opts.MaxMessageSize > 0 && header.Length > c.opts.MaxMessageSize {
					c.opts.handleError(fmt.Errorf("message too large: %d > %d", header.Length, c.opts.MaxMessageSize))
					return
				}

				// 从缓冲区池获取内存，减少分配
				payload := c.bufPool.Get().([]byte)[:header.Length]
				// 读取完整的消息体
				if _, err = io.ReadFull(c.raw, payload); err != nil {
					c.bufPool.Put(payload) // 出错时归还缓冲区
					c.opts.handleError(fmt.Errorf("read message read full failed: %v", err))
					return
				}

				// 根据消息类型触发相应的回调
				switch header.OpCode {
				case ws.OpText:
					c.opts.handleMessage(c, payload)
				case ws.OpBinary:
					c.opts.handleBinaryMessage(c, payload)
				}

				c.bufPool.Put(payload) // 处理完成后归还缓冲区
			}
		}
	}
}

// Close 关闭连接的外部方法
// 向done通道发送信号，通知读写协程退出
func (c *Conn) Close() {
	c.done <- struct{}{}
}

// close 连接内部的关闭方法
// 关闭通道和底层连接，释放资源
func (c *Conn) close() {
	// 关闭done通道和写入队列
	close(c.done)
	close(c.writeQueue)

	// 关闭底层网络连接
	if err := c.raw.Close(); err != nil {
		c.opts.handleError(fmt.Errorf("conn done failed: %v", err))
	}
}
