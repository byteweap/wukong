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

func (ss *Sessions) register(uid int64, s *melody.Session) {
	ss.data.Store(uid, s)
}

func (ss *Sessions) unregister(uid int64) {
	ss.data.Delete(uid)
}

func (ss *Sessions) get(uid int64) (*melody.Session, bool) {
	if session, ok := ss.data.Load(uid); ok {
		return session.(*melody.Session), true
	}
	return nil, false
}
