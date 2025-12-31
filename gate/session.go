package gate

import (
	"sync"
	"sync/atomic"

	"github.com/byteweap/wukong/plugin/network"
)

type Session struct {
	id   int64 // 会话唯一id
	uid  int64 // 用户id
	conn network.Conn
}

func newSession(id int64, uid int64, conn network.Conn) *Session {
	return &Session{
		id:   id,
		uid:  uid,
		conn: conn,
	}
}

func (s *Session) ID() int64 {
	return s.id
}

func (s *Session) UID() int64 {
	return s.uid
}

type SessionManager struct {
	nextID   int64
	mux      sync.RWMutex
	sessions map[int64]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		nextID:   0,
		sessions: make(map[int64]*Session),
	}
}

func (sm *SessionManager) Add(uid int64, conn network.Conn) {

	id := atomic.AddInt64(&sm.nextID, 1)

	sm.mux.Lock()
	defer sm.mux.Unlock()

	sm.sessions[id] = newSession(id, uid, conn)
}

func (sm *SessionManager) Get(uid int64) *Session {

	sm.mux.RLock()
	defer sm.mux.RUnlock()

	return sm.sessions[uid]
}
