# Logrus Log

本包实现 `component/log` 的 logrus 适配器。

**特性**
- 使用 `logrus.Logger` 作为底层
- 按 `component/log` Level 映射到 logrus Level
- 支持 `logrus.FieldKeyMsg` 作为消息字段

**安装**

```bash
go get github.com/byteweap/wukong/contrib/log/logrus
```

**最小用法**

```go
package main

import (
  lgrs "github.com/byteweap/wukong/contrib/log/logrus"
  "github.com/byteweap/wukong/component/log"
  "github.com/sirupsen/logrus"
)

func main() {
  base := logrus.New()
  lg := lgrs.NewLogger(base)

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
}
```

**注意事项**
- 非字符串 key 会被忽略
- 当 level 高于 `logrus.Logger.Level` 时直接返回

