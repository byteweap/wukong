package event

import (
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

// ExitGame 退出游戏
func (h *EventHandler) ExitGame(ctx *mesh.Context, req *pb.ExitGameRequest) {
	log.Infof("ExitGame, ctx: %v, req: %v", ctx, req.String())
}
