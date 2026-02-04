module github.com/byteweap/wukong/contrib/log/logrus

go 1.25.5

require (
	github.com/byteweap/wukong v0.0.1
	github.com/sirupsen/logrus v1.9.4
)

require golang.org/x/sys v0.26.0 // indirect

replace github.com/byteweap/wukong => ../../..
