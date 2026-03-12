package room

type Room struct {
	id int
}

func New(id int) *Room {
	return &Room{
		id: id,
	}
}
