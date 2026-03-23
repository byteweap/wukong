# Broker

## 说明

Broker 是服务之间的高性能通信抽象，只提供通信原语（Pub/Sub/Request/Reply），不内置业务级限流、分级、合并等策略。

## 接口概览

核心接口位于 `component/broker`：

- `ID()` 实现标识
- `Pub` 发布消息
- `Sub` 订阅消息（可用 `SubQueue` 设置 Queue Group）
- `Request` 请求-响应
- `Reply` 回复请求
- `Close` 关闭连接

消息结构：

- `Message.Subject` 主题
- `Message.Reply` 回复地址
- `Message.Header` 消息头（可选）
- `Message.Data` 负载

## 最小用法

以下示例使用 `contrib/broker/nats`：

```go
package main

import (
  "context"

  "github.com/byteweap/wukong/component/broker"
  nb "github.com/byteweap/wukong/contrib/broker/nats"
)

func main() {
  b, err := nb.New(nb.WithURLs("nats://127.0.0.1:4222"))
  if err != nil {
    panic(err)
  }
  defer b.Close()

  _, _ = b.Sub(context.Background(), "wk.game.cmd", func(msg *broker.Message) {
    _ = msg.Data
  })

  _ = b.Pub(context.Background(), "wk.game.event", []byte("hello"))
}
```

## 实现（contrib）

- NATS Core: `github.com/byteweap/wukong/contrib/broker/nats`
