package server

import "github.com/byteweap/meta/examples/game/internal/service"

var _ service.IPlayerService = (*Server)(nil)

func (g *Server) NumPlayers() int {
	return g.playerSpace.NumPlayers()
}
