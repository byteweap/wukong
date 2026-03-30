package rpc

import "github.com/byteweap/meta/examples/game/internal/service"

type RpcHandler struct {
	gs service.IRoomService
}

func New(gs service.IRoomService) *RpcHandler {
	return &RpcHandler{gs: gs}
}
