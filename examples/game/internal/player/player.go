package player

type Player struct {
	id int64
}

func New(id int64) *Player {
	return &Player{
		id: id,
	}
}

func (p *Player) ID() int64 {
	return p.id
}
