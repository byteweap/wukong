package rpc

import (
	"github.com/byteweap/wukong/examples/game/internal/handler"
)

type RpcHandler struct {
	gs handler.IService
}

func New(gs handler.IService) *RpcHandler {
	return &RpcHandler{gs: gs}
}
