# NATS Broker (Core)

本包实现 `plugin/broker` 的 **NATS Core** 版本，目标是“薄、快、稳定”：

- 只封装 `Publish/Subscribe(Queue)/Request` 等通信原语
- 不内置业务级：限流/分级/合并/丢弃/房间 Actor 等（这些应在上层服务实现）

## 最小用法

```go
import (
  "context"
  "time"
  bnats "github.com/byteweap/wukong/contrib/broker/nats"
  "github.com/byteweap/wukong/plugin/broker"
)

func example() error {
  b, err := bnats.New(bnats.WithURLs("nats://127.0.0.1:4222"))
  if err != nil { return err }
  defer b.Close()

  // Queue Group：水平扩展的竞争消费者
  sub, err := b.Subscribe(context.Background(),
    "wk.game.cmd.route.001.v1",
    func(ctx context.Context, msg *broker.Message) {
      // 处理 msg.Data
      // 若是 request-reply，可向 msg.Reply 发布响应
    },
    broker.WithQueue("game-route-001"),
  )
  if err != nil { return err }
  defer sub.Unsubscribe()

  // request-reply
  ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
  defer cancel()
  resp, err := b.Request(ctx, "wk.player.qry.get.v1", []byte("hello"))
  _ = resp
  return err
}
```


