package internal

import (
	"github.com/byteweap/wukong/examples/game/internal/handler/event"
	"github.com/byteweap/wukong/examples/game/internal/handler/rpc"
	"github.com/byteweap/wukong/examples/game/internal/server"
	"github.com/byteweap/wukong/server/mesh"
)

func New() *server.Server {
	g := server.New()

	e := event.New(g)
	g.Route(1, 1, mesh.Wrap(e.EnterGame))
	g.Route(2, 1, mesh.Wrap(e.GameExit))

	r := rpc.New(g)
	g.RequestRoute("findRoom", "v1", mesh.WrapRequest(r.FindRoom))

	// todo

	return g
}
