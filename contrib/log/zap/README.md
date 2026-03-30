# Zap Log

本包实现 `component/log` 的 zap 适配器。

**特性**
- 使用 `zap.Logger` 作为底层
- 支持自定义消息字段 key
- kv 必须成对，异常时会写 warn

**安装**

```bash
go get github.com/byteweap/meta/contrib/log/zap
```

**最小用法**

```go
package main

import (
  "github.com/byteweap/meta/component/log"
  wzap "github.com/byteweap/meta/contrib/log/zap"
  "go.uber.org/zap"
)

func main() {
  base := zap.NewExample()
  lg := wzap.NewLogger(base)
  wzap.WithMessageKey("msg")(lg)

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
  _ = lg.Close()
}
```

**配置**
- `WithMessageKey` 设置消息字段 key（默认 `log.DefaultMessageKey`）

**注意事项**
- kv 必须成对，否则会记录 warning 并返回
- `Close()` 内部调用 `Sync()`

