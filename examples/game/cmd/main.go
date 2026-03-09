package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/contrib/registry/nacos"
	"github.com/byteweap/wukong/examples/game/internal"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	clientCfg := &constant.ClientConfig{
		NamespaceId:  "zhaobin",
		BeatInterval: 5000,
	}
	serverCfgs := []constant.ServerConfig{
		{
			IpAddr:      "10.80.1.67",
			Port:        18848,
			ContextPath: "/nacos",
		},
	}
	nc, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientCfg,
		ServerConfigs: serverCfgs,
	})
	if err != nil {
		panic(err)
	}
	r := nacos.New(nc)

	s := internal.New()
	id := rand.IntN(10)
	err = wukong.New(
		wukong.ID(fmt.Sprintf("game-%d", id)),
		wukong.Name("game"),
		wukong.Version("v1.0.0"),
		wukong.Metadata(map[string]string{"author": "Leo"}),
		wukong.Server(s),
		wukong.Registry(r),
	).Run()
	if err != nil {
		log.Info(err)
	}
}
