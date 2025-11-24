module github.com/byteweap/wukong/gate

go 1.25.4

require (
	github.com/byteweap/wukong v0.0.0-20251124062958-668f656c990b
	github.com/byteweap/wukong/contrib/logger/zerolog v0.0.0-20251124062958-668f656c990b
	github.com/byteweap/wukong/contrib/network/websocket v0.0.0-20251124062958-668f656c990b
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

replace github.com/byteweap/wukong => ../
