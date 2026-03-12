package seat

import "errors"

// Space 座位空间
type Space struct {
	seats []*Seat
}

func NewSpace(size int) *Space {
	return &Space{
		seats: make([]*Seat, 0, size),
	}
}

// GetSeat 获取座位
func (s *Space) GetSeat(id int) (*Seat, bool) {
	if id < len(s.seats) {
		return s.seats[id], true
	}
	return nil, false
}

// Register 注册座位
func (s *Space) Register(seat *Seat) error {
	id := seat.id
	_, ok := s.GetSeat(id)
	if ok {
		return errors.New("seat already exists")
	}
	s.seats[id] = seat
	return nil
}

// Unregister 注销座位
func (s *Space) Unregister(id int) {
	if id >= len(s.seats) {
		return
	}
	s.seats[id] = nil
}
