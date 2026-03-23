# Tencent Log

本包实现 `component/log` 的腾讯云 CLS 日志适配器。

**特性**
- 基于 CLS AsyncProducer 异步发送
- 自动写入 `level` 字段
- kv 自动转字符串（结构体走 JSON）

**安装**

```bash
go get github.com/byteweap/wukong/contrib/log/tencent
```

**最小用法**

```go
package main

import (
  "github.com/byteweap/wukong/component/log"
  tlog "github.com/byteweap/wukong/contrib/log/tencent"
)

func main() {
  lg, _ := tlog.NewLogger(
    tlog.WithEndpoint("ap-shanghai.cls.tencentcs.com"),
    tlog.WithTopicID("topic-id"),
    tlog.WithAccessKey("ak"),
    tlog.WithAccessSecret("sk"),
  )
  defer lg.Close()

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
}
```

**配置**
- `WithEndpoint` 设置 CLS 访问地址
- `WithTopicID` 设置 Topic ID
- `WithAccessKey` 设置 AccessKey
- `WithAccessSecret` 设置 AccessSecret

**注意事项**
- 可通过 `GetProducer()` 获取底层 Producer

