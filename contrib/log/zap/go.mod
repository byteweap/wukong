module github.com/byteweap/wukong/contrib/log/zap

go 1.25.5

require (
	github.com/byteweap/wukong v0.0.0
	go.uber.org/zap v1.27.1
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/byteweap/wukong => ../../..
