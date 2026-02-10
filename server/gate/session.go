package gate

import (
	"sync"
)

// Session 会话
type Session struct {
	uid int64
}

func (s *Session) UID() int64 {
	return s.uid
}

// Sessions 管理所有会话
type Sessions struct {
	ss sync.Map
}

func newSessions() *Sessions {
	return &Sessions{ss: sync.Map{}}
}

func (s *Sessions) Add(session *Session) {
	s.ss.Store(session.uid, session)
}

func (s *Sessions) Get(uid int64) (*Session, bool) {
	if session, ok := s.ss.Load(uid); ok {
		return session.(*Session), true
	}
	return nil, false
}
