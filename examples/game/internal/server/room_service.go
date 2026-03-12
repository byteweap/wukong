package server

import "github.com/byteweap/wukong/examples/game/internal/service"

var _ service.IRoomService = (*Server)(nil)

func (g *Server) NumRooms() int {
	return g.roomSpace.NumRooms()
}
