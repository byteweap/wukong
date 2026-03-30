package event

import (
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/examples/game/internal/pb"
	"github.com/byteweap/meta/server/mesh"
)

// ExitGame 退出游戏
func (h *EventHandler) ExitGame(ctx *mesh.Context, req *pb.ExitGameRequest) {
	log.Infof("ExitGame, ctx: %v, req: %v", ctx, req.String())
}
