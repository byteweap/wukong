module github.com/byteweap/wukong/contrib/locator/redis

go 1.25.3

require (
	github.com/byteweap/wukong v0.0.0-20251121092935-0821f549de54
	github.com/redis/go-redis/v9 v9.17.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace github.com/byteweap/wukong => ../../../
