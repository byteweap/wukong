package main

import (
	"github.com/byteweap/wukong/examples/game/internal/game"
	"github.com/byteweap/wukong/examples/game/internal/handler/event"
	"github.com/byteweap/wukong/examples/game/internal/handler/rpc"
	"github.com/byteweap/wukong/server/mesh"
)

func main() {
	g := game.New()
	eventHandler := event.New(g)
	rpcHandler := rpc.New(g)

	g.Route(1, 1, mesh.Wrap(eventHandler.EnterGame))
	g.Route(2, 1, mesh.Wrap(eventHandler.GameExit))
	g.RequestRoute("findRoom", "v1", mesh.WrapRequest(rpcHandler.FindRoom))
}
