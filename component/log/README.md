# Logger

## Usage

### Structured logging

```go
logger := log.NewStdLogger(os.Stdout)
// 字段与 Valuer
logger = log.With(logger,
    "service.name", "helloworld",
    "service.version", "v1.0.0",
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
)
logger.Log(log.LevelInfo, "key", "value")

// 辅助器
helper := log.NewHelper(logger)
helper.Log(log.LevelInfo, "key", "value")
helper.Info("info message")
helper.Infof("info %s", "message")
helper.Infow("key", "value")

// 过滤器
log := log.NewHelper(log.NewFilter(logger,
	log.FilterLevel(log.LevelInfo),
	log.FilterKey("foo"),
	log.FilterValue("bar"),
	log.FilterFunc(customFilter),
))
log.Debug("debug log")
log.Info("info log")
log.Warn("warn log")
log.Error("warn log")
```

## Third party log library

### zap

```shell
go get -u github.com/byteweap/wukong/contrib/log/zap/v2
```

### logrus

```shell
go get -u github.com/byteweap/wukong/contrib/log/logrus
```

### fluent

```shell
go get -u github.com/byteweap/wukong/contrib/log/fluent
```

### aliyun

```shell
go get -u github.com/byteweap/wukong/contrib/log/aliyun
```
