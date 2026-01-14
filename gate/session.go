package gate

import (
	"sync"
	"sync/atomic"

	"github.com/byteweap/wukong/component/network"
)

// Session 会话，表示一个用户连接
type Session struct {
	id   int64        // 会话唯一ID
	uid  int64        // 用户ID
	conn network.Conn // 网络连接
}

// newSession 创建新会话
func newSession(id int64, uid int64, conn network.Conn) *Session {
	return &Session{
		id:   id,
		uid:  uid,
		conn: conn,
	}
}

// ID 返回会话ID
func (s *Session) ID() int64 {
	return s.id
}

// UID 返回用户ID
func (s *Session) UID() int64 {
	return s.uid
}

// SessionManager 会话管理器，管理所有用户会话
type SessionManager struct {
	nextID   int64              // 下一个会话ID
	mux      sync.RWMutex       // 读写锁
	sessions map[int64]*Session // 会话映射表，key为uid
}

// NewSessionManager 创建会话管理器
func NewSessionManager() *SessionManager {
	return &SessionManager{
		nextID:   0,
		sessions: make(map[int64]*Session),
	}
}

// Add 添加会话
func (sm *SessionManager) Add(uid int64, conn network.Conn) {
	id := atomic.AddInt64(&sm.nextID, 1)

	sm.mux.Lock()
	defer sm.mux.Unlock()

	sm.sessions[uid] = newSession(id, uid, conn)
}

// Get 获取指定用户的会话
func (sm *SessionManager) Get(uid int64) *Session {
	sm.mux.RLock()
	defer sm.mux.RUnlock()

	return sm.sessions[uid]
}

// Remove 移除指定用户的会话
func (sm *SessionManager) Remove(uid int64) {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	delete(sm.sessions, uid)
}

// Close 关闭所有会话
func (sm *SessionManager) Close() error {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	for _, session := range sm.sessions {
		session.conn.Close()
	}

	return nil
}
