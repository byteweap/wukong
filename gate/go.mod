module github.com/byteweap/wukong/gate

go 1.25.4

require (
	github.com/byteweap/wukong v0.0.0-20251125012552-766f6f541a96
	github.com/byteweap/wukong/contrib/locator/redis v0.0.0-20251231091904-f21ac936137d
	github.com/byteweap/wukong/contrib/logger/zerolog v0.0.0-20251125012552-766f6f541a96
	github.com/byteweap/wukong/contrib/network/websocket v0.0.0-20251125012552-766f6f541a96
	github.com/redis/go-redis/v9 v9.17.0
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/byteweap/wukong => ../
