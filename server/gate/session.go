package gate

import (
	"sync"

	"github.com/olahol/melody"
)

type Event int8

const (
	EventConnect    Event = iota // 建立链接
	EventDisconnect              // 断开链接
	EventReconnect               // 重连
)

// Sessions 管理所有会话
type Sessions struct {
	data sync.Map
}

func newSessions() *Sessions {
	return &Sessions{data: sync.Map{}}
}

// register 注册会话
func (ss *Sessions) register(uid int64, s *melody.Session) {
	ss.data.Store(uid, s)
}

// unregister 注销会话
func (ss *Sessions) unregister(uid int64) {
	ss.data.Delete(uid)
}

// get 获取会话
func (ss *Sessions) get(uid int64) (*melody.Session, bool) {
	if session, ok := ss.data.Load(uid); ok {
		return session.(*melody.Session), true
	}
	return nil, false
}
