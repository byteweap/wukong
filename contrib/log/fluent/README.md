# Fluent Log

本包实现 `component/log` 的 fluentd 日志适配器。

**特性**
- 支持 `tcp://host:port` 与 `unix://path` 两种地址
- 可配置超时、重试与缓冲
- kv 结构自动转为 map 发送

**安装**

```bash
go get github.com/byteweap/wukong/contrib/log/fluent
```

**最小用法**

```go
package main

import (
  "github.com/byteweap/wukong/component/log"
  flog "github.com/byteweap/wukong/contrib/log/fluent"
)

func main() {
  lg, _ := flog.NewLogger("tcp://127.0.0.1:24224")
  defer lg.Close()

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
}
```

**配置**
- `WithTimeout` 设置超时
- `WithWriteTimeout` 设置写超时
- `WithBufferLimit` 设置缓冲大小
- `WithRetryWait` 设置重试间隔
- `WithMaxRetry` 设置最大重试次数
- `WithMaxRetryWait` 设置最大重试等待
- `WithTagPrefix` 设置 tag 前缀
- `WithAsync` 启用异步发送
- `WithForceStopAsyncSend` 强制停止异步发送

**注意事项**
- kv 为奇数时会自动追加 `KEYVALS UNPAIRED`
- `addr` 必须包含 scheme：`tcp` 或 `unix`

