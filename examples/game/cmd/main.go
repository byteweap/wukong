package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/byteweap/meta"
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/contrib/registry/nacos"
	"github.com/byteweap/meta/examples/game/internal"
)

func newNamingClient() naming_client.INamingClient {
	clientCfg := constant.ClientConfig{
		//NamespaceId:  "zhaobin",
		//BeatInterval: 5000,
	}
	serverCfgs := []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			Port:        18848,
			ContextPath: "/nacos",
		},
	}
	nc, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientCfg,
			ServerConfigs: serverCfgs,
		},
	)
	if err != nil {
		panic(err)
	}
	return nc
}

func main() {

	// 1. nacos注册中心
	nc := newNamingClient()
	registry := nacos.New(nc)

	// 2. server [mesh]
	s, cleanup, err := internal.New()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	err = meta.New(
		meta.ID(fmt.Sprintf("game-%d", rand.IntN(100))),
		meta.Name("game"),
		meta.Version("v1.0.0"),
		meta.Metadata(map[string]string{"author": "Leo"}),
		meta.Server(s),
		meta.Registry(registry),
	).Run()
	if err != nil {
		log.Info(err)
	}
}
