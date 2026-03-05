package event

import (
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

func (h *EventHandler) GameExit(ctx *mesh.Context, req *pb.ExitGameRequest) {
	log.Infof("GameExit, ctx: %v, req: %v", ctx, req.String())
}
