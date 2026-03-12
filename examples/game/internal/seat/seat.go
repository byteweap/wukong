package seat

import "github.com/byteweap/wukong/examples/game/internal/player"

// Seat 座位
type Seat struct {
	id int
	p  *player.Player
}

func New(id int) *Seat {
	return &Seat{
		id: id,
	}
}

func (s *Seat) ID() int {
	return s.id
}

func (s *Seat) Player() *player.Player {
	return s.p
}
