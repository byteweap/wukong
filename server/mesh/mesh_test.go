package mesh_test

import (
	"testing"

	"github.com/byteweap/wukong/server/mesh"
)

type Params struct {
	Name string
}

// EnterGame 模拟业务处理函数
// pub-sub
func EnterGame(ctx *mesh.Context, req *Params) {
}

// FindUser 模拟业务处理函数
// request-reply
func FindUser(ctx *mesh.RequestContext, req *Params) {
	return
}

// TestMesh 验证 Route 注册基本可用
func TestMesh(t *testing.T) {
	app := mesh.New()
	app.Route(1, 1, EnterGame)
	app.RequestRoute("findUser", "v1", FindUser)
}
