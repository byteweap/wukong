package rpc

import "github.com/byteweap/wukong/examples/game/internal/game"

type RpcHandler struct {
	g *game.Game
}

func New(g *game.Game) *RpcHandler {
	return &RpcHandler{g: g}
}
