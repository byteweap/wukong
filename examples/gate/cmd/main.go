package main

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	goredis "github.com/redis/go-redis/v9"

	"github.com/byteweap/wukong"
	"github.com/byteweap/wukong/component/log"
	"github.com/byteweap/wukong/contrib/broker/nats"
	"github.com/byteweap/wukong/contrib/locator/redis"
	"github.com/byteweap/wukong/contrib/registry/nacos"
	"github.com/byteweap/wukong/server/gate"
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
	loc := redis.New(goredis.UniversalOptions{}, "wukong")
	defer loc.Close()

	// 3. broker
	broker, err := nats.New()
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	err = wukong.New(
		wukong.ID("gate-1"),
		wukong.Name("gate"),
		wukong.Version("v1.0.0"),
		wukong.Metadata(map[string]string{"author": "Leo"}),
		wukong.Server(gate.New(
			gate.Locator(loc),
			gate.Discovery(reg),
			gate.Broker(broker),
		)),
		wukong.Registry(reg),
	).Run()
	if err != nil {
		log.Info(err)
	}
}
