# Zerolog Log

本包实现 `component/log` 的 zerolog 适配器。

**特性**
- 使用 `zerolog.Logger` 作为底层
- level 映射到对应事件
- 仅接收 string key

**安装**

```bash
go get github.com/byteweap/wukong/contrib/log/zerolog
```

**最小用法**

```go
package main

import (
  "os"

  "github.com/byteweap/wukong/component/log"
  zlog "github.com/byteweap/wukong/contrib/log/zerolog"
  "github.com/rs/zerolog"
)

func main() {
  base := zerolog.New(os.Stdout)
  lg := zlog.NewLogger(&base)

  _ = lg.Log(log.LevelInfo, "msg", "hello", "uid", 1001)
}
```

**注意事项**
- 非字符串 key 会被忽略
- kv 为奇数时会自动追加空值
