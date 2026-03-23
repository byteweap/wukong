# Logger

## 说明

log 组件提供统一的结构化日志接口，包含日志级别、字段、辅助器与过滤器等能力。

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

## Contrib 实现

### zap

```shell
go get -u github.com/byteweap/wukong/contrib/log/zap
```

### logrus

```shell
go get -u github.com/byteweap/wukong/contrib/log/logrus
```

### zerolog

```shell
go get -u github.com/byteweap/wukong/contrib/log/zerolog
```

### fluent

```shell
go get -u github.com/byteweap/wukong/contrib/log/fluent
```

### aliyun

```shell
go get -u github.com/byteweap/wukong/contrib/log/aliyun
```

### tencent

```shell
go get -u github.com/byteweap/wukong/contrib/log/tencent
```
