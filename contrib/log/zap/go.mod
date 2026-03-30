module github.com/byteweap/meta/contrib/log/zap

go 1.26.1

require (
	github.com/byteweap/meta v0.0.1
	go.uber.org/zap v1.27.1
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/byteweap/meta => ../../..
