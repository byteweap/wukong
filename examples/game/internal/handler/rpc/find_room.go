package rpc

import (
	"net/http"

	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

func (h *RpcHandler) FindRoom(ctx *mesh.RpcContext, req *pb.FindRoomRequest) ([]byte, string, int) {
	//pb.
	return nil, "", http.StatusOK
}
