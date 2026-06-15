package gate

import (
	"sync"

	"github.com/olahol/melody"
)

// Sessions 管理所有会话
type Sessions struct {
	data sync.Map
}

func newSessions() *Sessions {
	return &Sessions{data: sync.Map{}}
}

// replace 注册新会话并返回旧会话
func (ss *Sessions) replace(uid int64, s *melody.Session) (*melody.Session, bool) {
	old, loaded := ss.data.Swap(uid, s)
	if !loaded {
		return nil, false
	}
	return old.(*melody.Session), true
}

// unregisterIfSame 仅当当前会话匹配时注销
func (ss *Sessions) unregisterIfSame(uid int64, s *melody.Session) bool {
	return ss.data.CompareAndDelete(uid, s)
}

// get 获取会话
func (ss *Sessions) get(uid int64) (*melody.Session, bool) {
	if session, ok := ss.data.Load(uid); ok {
		return session.(*melody.Session), true
	}
	return nil, false
}
