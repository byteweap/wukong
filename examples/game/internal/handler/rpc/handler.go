package rpc

import (
	"github.com/byteweap/wukong/examples/game/internal/handler"
)

type RpcHandler struct {
	gs route.GameService
}

func New(gs route.GameService) *RpcHandler {
	return &RpcHandler{gs: gs}
}
