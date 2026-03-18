package internal

import (
	goredis "github.com/redis/go-redis/v9"

	"github.com/byteweap/wukong/contrib/broker/nats"
	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/examples/game/internal/handler/event"
	"github.com/byteweap/wukong/examples/game/internal/handler/rpc"
	"github.com/byteweap/wukong/examples/game/internal/server"
	"github.com/byteweap/wukong/server/mesh"
)

func New() (*server.Server, func(), error) {

	loc := redis.New(goredis.UniversalOptions{
		Addrs: []string{"localhost:6379"},
	}, "game")

	bro, err := nats.New(
		nats.URLs("nats://localhost:4222"),
	)
	if err != nil {
		return nil, nil, err
	}

	g := server.New(
		mesh.Broker(bro),
		mesh.Locator(loc),
	)

	e := event.New(g)
	g.Route(1, 1, mesh.Wrap(e.EnterGame))
	g.Route(2, 1, mesh.Wrap(e.ExitGame))

	r := rpc.New(g)
	g.RpcRoute("hello", "v1", mesh.WrapRpc(r.Hello))

	return g, func() {
		_ = loc.Close()
		_ = bro.Close()
	}, nil
}
