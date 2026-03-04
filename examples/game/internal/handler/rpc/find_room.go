package rpc

import (
	"net/http"

	"github.com/byteweap/wukong/server/mesh"
)

type FindRoomParams struct {
}

func (h *RpcHandler) FindRoom(ctx *mesh.RequestContext, req *FindRoomParams) ([]byte, string, int) {
	//todo
	return nil, "", http.StatusOK
}
