# Locator

## 说明

Locator 用于跟踪玩家在各服务节点上的位置，提供绑定、解绑与查询能力。

## 接口概览

核心接口位于 `component/locator`：

- `ID()` 实现标识
- `AllNodes` 查询用户在所有服务上的节点
- `Node` 查询用户在某服务上的节点
- `Bind` 绑定用户到某服务节点
- `UnBind` 解绑用户与某服务节点
- `Close` 关闭定位器

## 最小用法

以下示例使用 `contrib/locator/redis`：

```go
package main

import (
  "context"

  rloc "github.com/byteweap/wukong/contrib/locator/redis"
  "github.com/redis/go-redis/v9"
)

func main() {
  opts := redis.UniversalOptions{Addrs: []string{"127.0.0.1:6379"}}
  loc := rloc.New(opts, "wk")
  defer loc.Close()

  ctx := context.Background()
  _ = loc.Bind(ctx, 1001, "gate", "node-1")
  _, _ = loc.Node(ctx, 1001, "gate")
  _ = loc.UnBind(ctx, 1001, "gate", "node-1")
}
```

## 实现（contrib）

- Redis: `github.com/byteweap/wukong/contrib/locator/redis`
