package event

import (
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/examples/game/internal/pb"
	"github.com/byteweap/meta/server/mesh"
)

// EnterGame 进入游戏
func (h *EventHandler) EnterGame(ctx *mesh.Context, req *pb.EnterGameRequest) {
	log.Infof("EnterGame, ctx: %v, req: %v", ctx, req.String())
	ctx.OkResp()
}
