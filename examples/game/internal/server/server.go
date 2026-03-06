package server

import (
	"sync"

	"github.com/byteweap/wukong/examples/game/internal/room"
	"github.com/byteweap/wukong/server/mesh"
)

type Server struct {
	*mesh.Mesh

	mu    sync.RWMutex
	rooms map[int]*room.Room
}

func New(opts ...mesh.Option) *Server {
	m := mesh.New(opts...)
	return &Server{
		Mesh:  m,
		rooms: make(map[int]*room.Room),
	}
}

func (g *Server) NumRooms() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.rooms)
}
