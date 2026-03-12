package room

import (
	"sync"
	"sync/atomic"
)

// Space 房间空间
type Space struct {
	num   atomic.Uint32 // 房间数量
	rooms sync.Map      // 房间集合 key: 房间ID value: *Room
}

func NewSpace() *Space {
	return &Space{}
}

// GetRoom 获取房间
func (s *Space) GetRoom(id int) (*Room, bool) {
	v, ok := s.rooms.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*Room), true
}

// Register 注册房间
func (s *Space) Register(r *Room) {
	s.num.Add(1)
	s.rooms.Store(r.id, r)
}

// Unregister 注销房间
func (s *Space) Unregister(id int) {
	s.num.Add(-1)
	s.rooms.Delete(id)
}

// NumRooms 获取房间数量
func (s *Space) NumRooms() int {
	return int(s.num.Load())
}
