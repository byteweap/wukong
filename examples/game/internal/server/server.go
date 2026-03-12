package server

import (
	"github.com/byteweap/wukong/examples/game/internal/player"
	"github.com/byteweap/wukong/examples/game/internal/room"
	"github.com/byteweap/wukong/server/mesh"
)

// Server 核心服务
type Server struct {
	*mesh.Mesh

	roomSpace   *room.Space   // 房间空间
	playerSpace *player.Space // 玩家空间
}

func New(opts ...mesh.Option) *Server {
	m := mesh.New(opts...)
	return &Server{
		Mesh:        m,
		roomSpace:   room.NewSpace(),
		playerSpace: player.NewSpace(),
	}
}
