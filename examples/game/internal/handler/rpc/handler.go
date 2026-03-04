package rpc

import (
	"github.com/byteweap/wukong/examples/game/internal/handler"
)

type RpcHandler struct {
	gs handler.GameService
}

func New(gs handler.GameService) *RpcHandler {
	return &RpcHandler{gs: gs}
}
