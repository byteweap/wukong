package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	goredis "github.com/redis/go-redis/v9"

	"github.com/byteweap/meta"
	"github.com/byteweap/meta/component/log"
	"github.com/byteweap/meta/component/selector"
	"github.com/byteweap/meta/contrib/broker/nats"
	"github.com/byteweap/meta/contrib/locator/redis"
	"github.com/byteweap/meta/contrib/registry/nacos"
	"github.com/byteweap/meta/contrib/selector/wrr"
	"github.com/byteweap/meta/server/gate"
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
	reg := nacos.New(newNamingClient())

	// 2. 定位器
	loc := redis.New(goredis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
	}, "meta")
	defer loc.Close()

	// 3. broker
	broker, err := nats.New()
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	err = meta.New(
		meta.ID(fmt.Sprintf("gate-%d", rand.IntN(100))),
		meta.Name("gate"),
		meta.Version("v1.0.0"),
		meta.Metadata(map[string]string{"author": "Leo"}),
		meta.Server(gate.New(
			gate.Addr(fmt.Sprintf(":%d", rand.IntN(1000)+8000)),
			gate.Locator(loc),
			gate.Discovery(reg),
			gate.Broker(broker),
			gate.SelectorFunc(func() selector.Selector {
				return wrr.New()
			}),
		)),
		meta.Registry(reg),
	).Run()
	if err != nil {
		log.Errorf("app run error: %v", err)
	}
}
