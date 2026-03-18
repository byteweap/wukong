package rpc

import (
	"net/http"

	"github.com/byteweap/wukong/examples/game/internal/pb"
	"github.com/byteweap/wukong/server/mesh"
)

// Hello RPC示例接口
func (h *RpcHandler) Hello(ctx *mesh.RpcContext, req *pb.FindRoomRequest) ([]byte, string, int) {
	return []byte("Hello RPC"), "", http.StatusOK
}
