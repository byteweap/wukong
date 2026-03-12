package server

import "github.com/byteweap/wukong/examples/game/internal/handler"

var _ handler.IService = (*Server)(nil)

func (g *Server) NumRooms() int {
	return g.roomSpace.NumRooms()
}
