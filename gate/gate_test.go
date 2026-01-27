package gate_test

//
//import (
//	"testing"
//	"time"
//
//	"github.com/byteweap/wukong/contrib/broker/nats"
//	"github.com/byteweap/wukong/contrib/locator/redis"
//	"github.com/byteweap/wukong/contrib/network/websocket"
//	"github.com/byteweap/wukong/contrib/registry/nacos"
//	"github.com/byteweap/wukong/gate"
//	"github.com/nacos-group/nacos-sdk-go/clients"
//	"github.com/nacos-group/nacos-sdk-go/common/constant"
//	"github.com/nacos-group/nacos-sdk-go/vo"
//	redis2 "github.com/redis/go-redis/v9"
//	"gopkg.in/natefinch/lumberjack.v2"
//)
//
//var testServerConfig = []constant.ServerConfig{
//	*constant.NewServerConfig("127.0.0.1", 18848),
//}
//
//func TestGate(t *testing.T) {
//	loc := redis.New(redis2.UniversalOptions{
//		Addrs: []string{":6379"},
//	}, "data:%d", "gate", "game")
//
//	bro, err := nats.New()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cc := constant.ClientConfig{
//		NamespaceId:         "public", // 命名空间 id
//		TimeoutMs:           5000,
//		NotLoadCacheAtStart: true,
//		LogDir:              "/tmp/nacos/log",
//		CacheDir:            "/tmp/nacos/cache",
//		LogRollingConfig:    &lumberjack.Logger{},
//		LogLevel:            "debug",
//	}
//
//	// 更稳妥的方式创建 naming client
//	client, err := clients.NewNamingClient(
//		vo.NacosClientParam{
//			ClientConfig:  &cc,
//			ServerConfigs: testServerConfig,
//		},
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	reg := nacos.New(client)
//
//	g := gate.New(
//		gate.ID("aaaa"),
//		gate.Name("ivy-gateway"),
//		gate.Version("v0.0.1"),
//		gate.Metadata(map[string]string{"abc": "123"}),
//		gate.NetServer(websocket.NewServer(
//			websocket.Addr(":8088"),
//		)),
//		gate.Locator(loc),
//		gate.Broker(bro),
//		gate.Registry(reg, time.Second*5),
//	)
//
//	if err = g.Run(); err != nil {
//		t.Fatal(err)
//	}
//	if err = g.Stop(); err != nil {
//		t.Fatal(err)
//	}
//}
