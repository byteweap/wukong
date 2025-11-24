package websocket

import (
	"math"
	"net"
	"sync"
	"sync/atomic"
)

// hub 管理所有WebSocket连接的连接管理器
type hub struct {
	opts       *Options    // 服务器配置选项
	nextConnID int64       // 下一个连接ID
	open       atomic.Bool // 是否打开（原子操作）

	mux     sync.RWMutex         // 读写锁，保护conns映射
	connNum atomic.Int64         // 当前连接数（原子操作）
	conns   map[*wsConn]struct{} // 所有活跃连接映射
}

// newHub 创建新的连接管理器实例
// opts: 服务器配置选项
// 返回: 初始化好的hub实例
func newHub(opts *Options) *hub {
	h := &hub{
		opts:       opts,
		nextConnID: 0,
		conns:      make(map[*wsConn]struct{}),
	}
	h.open.Store(true) // 初始状态为打开
	h.connNum.Store(0) // 初始连接数为0
	return h
}

// closed 判断hub是否已关闭
func (h *hub) closed() bool {
	return !h.open.Load()
}

// shutdown 优雅关闭hub
func (h *hub) shutdown() error {
	if h.closed() {
		return ErrHubClosed
	}
	h.open.Store(false)

	// 关闭所有活跃连接
	h.mux.RLock()
	for conn := range h.conns {
		conn.Close()
	}
	h.mux.RUnlock()

	// 清空连接映射
	clear(h.conns)

	return nil
}

// nextId 生成下一个连接ID
func (h *hub) nextId() int64 {
	// 原子递增连接ID
	atomic.AddInt64(&h.nextConnID, 1)

	// 处理ID溢出，重置为1
	if atomic.LoadInt64(&h.nextConnID) == math.MaxInt64 {
		atomic.StoreInt64(&h.nextConnID, 1)
	}
	return atomic.LoadInt64(&h.nextConnID)
}

// allocate 分配的WebSocket连接
func (h *hub) allocate(netConn net.Conn) error {
	if h.closed() {
		return ErrHubClosed
	}
	if h.connNum.Load() >= int64(h.opts.MaxConnections) {
		return ErrMaxConns
	}
	id := h.nextId()
	conn := newConn(id, netConn, h.opts)

	// 添加到连接映射，增加计数
	h.mux.Lock()
	h.conns[conn] = struct{}{}
	h.mux.Unlock()
	h.connNum.Add(1)
	h.opts.handleConnect(conn)

	go conn.writePump() // 启动写协程
	conn.readPump()     // 读协程

	// 连接关闭后，从映射中移除, 减少计数
	h.mux.Lock()
	delete(h.conns, conn)
	h.mux.Unlock()
	h.connNum.Add(-1)

	// 触发连接断开回调
	h.opts.handleDisconnect(conn)

	return nil
}
