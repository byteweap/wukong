package event

import (
	"github.com/byteweap/meta/examples/game/internal/pb"
	"github.com/byteweap/meta/server/mesh"
)

// Hello event 示例接口
func (h *EventHandler) Hello(ctx *mesh.Context, req *pb.ExitGameRequest) {
	ctx.OkResp()
}
