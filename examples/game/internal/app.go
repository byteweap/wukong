package internal

import (
	"github.com/byteweap/wukong/contrib/broker/nats"
	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/examples/game/internal/handler/event"
	"github.com/byteweap/wukong/examples/game/internal/handler/rpc"
	"github.com/byteweap/wukong/examples/game/internal/server"
	"github.com/byteweap/wukong/server/mesh"
	goredis "github.com/redis/go-redis/v9"
)

func New() *server.Server {

	loc := redis.New(goredis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	}, "game")

	bro, err := nats.New(
		nats.URLs("nats://localhost:4222"),
	)
	if err != nil {
		panic(err)
	}

	g := server.New(
		mesh.Broker(bro),
		mesh.Locator(loc),
	)

	e := event.New(g)
	g.Route(1, 1, mesh.Wrap(e.EnterGame))
	g.Route(2, 1, mesh.Wrap(e.GameExit))

	r := rpc.New(g)
	g.RpcRoute("findRoom", "v1", mesh.WrapRpc(r.FindRoom))

	// todo

	return g
}
