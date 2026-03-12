package player

import (
	"sync"
	"sync/atomic"
)

// Space 玩家空间
type Space struct {
	num     atomic.Uint32 // 玩家数量
	players sync.Map      // 玩家集合 key: 玩家ID value: *Player
}

func NewSpace() *Space {
	return &Space{}
}

// GetPlayer 获取玩家
func (s *Space) GetPlayer(id int64) (*Player, bool) {
	p, ok := s.players.Load(id)
	if !ok {
		return nil, false
	}
	return p.(*Player), true
}

// Register 注册玩家
func (s *Space) Register(p *Player) {
	s.players.Store(p.ID(), p)
	s.num.Add(1)
}

// Unregister 注销玩家
func (s *Space) Unregister(id int64) {
	s.players.Delete(id)
	s.num.Add(-1)
}

// NumPlayers 获取玩家数量
func (s *Space) NumPlayers() int {
	return int(s.num.Load())
}
