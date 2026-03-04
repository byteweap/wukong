package game

import (
	"sync"

	"github.com/byteweap/wukong/examples/game/internal/room"
	"github.com/byteweap/wukong/server/mesh"
)

type Game struct {
	*mesh.Mesh

	mu    sync.RWMutex
	rooms map[int]*room.Room
}

func New() *Game {
	m := mesh.New()
	return &Game{
		Mesh:  m,
		rooms: make(map[int]*room.Room),
	}
}

func (g *Game) NumRooms() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.rooms)
}
