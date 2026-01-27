module github.com/byteweap/wukong/gate

go 1.25.5

require (
	github.com/byteweap/wukong v0.0.1
	github.com/google/uuid v1.6.0
)

replace (
	github.com/byteweap/wukong => ..
	github.com/byteweap/wukong/contrib/locator/redis => ../contrib/locator/redis
	github.com/byteweap/wukong/contrib/network/websocket => ../contrib/network/websocket
	github.com/byteweap/wukong/contrib/registry/nacos => ../contrib/registry/nacos
)
