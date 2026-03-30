# Redis Locator

本包实现 `component/locator` 的 Redis 版本，使用 Hash 保存用户在不同服务上的节点归属。

**特性**
- 基于 Redis Hash 存储 `service -> node`
- 支持统一前缀与默认键规则
- 提供完整的绑定/解绑/查询能力

**安装**

```bash
go get github.com/byteweap/meta/contrib/locator/redis
```

**最小用法**

```go
package main

import (
  "context"

  rloc "github.com/byteweap/meta/contrib/locator/redis"
  "github.com/redis/go-redis/v9"
)

func main() {
  opts := redis.UniversalOptions{
    Addrs: []string{"127.0.0.1:6379"},
  }

  loc := rloc.New(opts, "wk")
  defer loc.Close()

  ctx := context.Background()
  _ = loc.Bind(ctx, 1001, "gate", "node-1")

  node, _ := loc.Node(ctx, 1001, "gate")
  _ = node

  _ = loc.UnBind(ctx, 1001, "gate", "node-1")
}
```

**配置**
- `New(opts, prefix)` 使用 `redis.UniversalOptions`
- `prefix` 为空时使用默认键规则

**注意事项**
- 键规则：`<prefix>:locator:<uid>`，当 `prefix` 为空时为 `locator:<uid>`
- `UnBind` 仅在当前节点匹配时才删除映射

