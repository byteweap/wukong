package event

import (
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

func (h *EventHandler) EnterGame(ctx *mesh.Context, req *pb.EnterGameRequest) {
	log.Infof("EnterGame, ctx: %v, req: %v", ctx, req.String())
	ctx.OkResp()
}
