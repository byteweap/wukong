module github.com/byteweap/meta/examples/game

go 1.26.1

require (
	github.com/byteweap/meta v0.0.1
	github.com/byteweap/meta/contrib/broker/nats v0.0.0-00010101000000-000000000000
	github.com/byteweap/meta/contrib/locator/redis v0.0.0-00010101000000-000000000000
	github.com/byteweap/meta/contrib/registry/nacos v0.0.0-00010101000000-000000000000
	github.com/nacos-group/nacos-sdk-go v1.1.6
	github.com/redis/go-redis/v9 v9.18.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1800 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats.go v1.48.0 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/byteweap/meta => ../..

replace github.com/byteweap/meta/contrib/broker/nats => ../../contrib/broker/nats

replace github.com/byteweap/meta/contrib/locator/redis => ../../contrib/locator/redis

replace github.com/byteweap/meta/contrib/registry/nacos => ../../contrib/registry/nacos
