package event

import (
	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

func (h *EventHandler) GameExit(ctx *mesh.Context, req *pb.ExitGameRequest) {
	// todo
}
