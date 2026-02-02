package gate

import "github.com/byteweap/wukong/component/network"

type Session struct {
	raw network.Conn
	uid int64
}

func newSession(raw network.Conn, uid int64) *Session {
	return &Session{raw: raw, uid: uid}
}

func (s *Session) Raw() network.Conn {
	return s.raw
}

func (s *Session) UID() int64 {
	return s.uid
}
