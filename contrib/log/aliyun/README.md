# Aliyun Log

本包实现 `component/log` 的阿里云 SLS 日志适配器。

**特性**
- 基于 SLS Producer 异步发送
- 自动写入 `level` 字段
- 将 kv 统一转为字符串（结构体走 JSON）

**安装**

```bash
go get github.com/byteweap/wukong/contrib/log/aliyun
```

**最小用法**

```go
package main

import (
  "github.com/byteweap/wukong/component/log"
  alog "github.com/byteweap/wukong/contrib/log/aliyun"
)

func main() {
  lg, _ := alog.NewAliyunLog(
    alog.WithEndpoint("cn-shanghai.log.aliyuncs.com"),
    alog.WithProject("projectName"),
    alog.WithLogstore("app"),
    alog.WithAccessKey("ak"),
    alog.WithAccessSecret("sk"),
  )
  defer lg.Close()

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
}
```

**配置**
- `WithEndpoint` 设置 SLS 访问地址
- `WithProject` 设置 Project
- `WithLogstore` 设置 Logstore
- `WithAccessKey` 设置 AccessKey
- `WithAccessSecret` 设置 AccessSecret

**注意事项**
- 默认 `project` 为 `projectName`，`logstore` 为 `app`
- 可通过 `GetProducer()` 获取底层 Producer

