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
func FindUser(ctx *mesh.RpcContext, req *Params) ([]byte, string, int) {
	return nil, "ok", 200
}

// TestMesh
func TestMesh(t *testing.T) {
	app := mesh.New()

	// 1. Gate 消息路由
	// 	1.1. 写法简单,但使用反射,有性能开销, 不推荐用于高频路由
	app.RouteX(1, 1, EnterGame)
	// 	1.2. 推荐写法, 避免反射调用, 推荐用于高频路由
	app.Route(2, 1, mesh.Wrap(EnterGame))

	// 2. Request 消息路由
	// 	2.1. 写法简单,但使用反射,有性能开销, 不推荐用于高频路由
	app.RequestRouteX("findUser", "v1", FindUser)
	// 	2.2. 推荐写法, 避免反射调用, 推荐用于高频路由
	app.RequestRoute("findUser1", "v1", mesh.WrapRpc(FindUser))

}
