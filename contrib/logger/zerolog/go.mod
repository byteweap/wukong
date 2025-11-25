module github.com/byteweap/wukong/contrib/logger/zerolog

go 1.25.4

require (
	github.com/byteweap/wukong v0.0.0-20251124080543-a142f476f77e
	github.com/rs/zerolog v1.34.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

replace github.com/byteweap/wukong => ../../../
