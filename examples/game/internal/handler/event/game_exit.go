package event

import (
	"github.com/byteweap/wukong/server/mesh"
)

type GameExitParams struct {
}

func (h *EventHandler) GameExit(ctx *mesh.Context, req *GameExitParams) {
	// todo
}
