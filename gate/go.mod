module github.com/byteweap/wukong/gate

go 1.25.5

require (
	github.com/byteweap/wukong v0.0.0
	github.com/byteweap/wukong/contrib/broker/nats v0.0.0-20260123061731-15d484da05a5
	github.com/byteweap/wukong/contrib/locator/redis v0.0.0
	github.com/byteweap/wukong/contrib/network/websocket v0.0.0
	github.com/byteweap/wukong/contrib/registry/nacos v0.0.0
	github.com/google/uuid v1.6.0
	github.com/nacos-group/nacos-sdk-go v1.1.6
	github.com/redis/go-redis/v9 v9.17.2
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.18 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/nats-io/nats.go v1.48.0 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
)

replace (
	github.com/byteweap/wukong => ..
	github.com/byteweap/wukong/contrib/locator/redis => ../contrib/locator/redis
	github.com/byteweap/wukong/contrib/network/websocket => ../contrib/network/websocket
	github.com/byteweap/wukong/contrib/registry/nacos => ../contrib/registry/nacos
)
