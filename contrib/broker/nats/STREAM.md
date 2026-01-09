# JetStream 扩展使用指南

## 设计理念

采用**渐进式引入**策略：
- **Core Broker**（已完成）：用于 90% 的实时链路（战斗指令、状态同步）
- **Stream Broker**（可选扩展）：用于 10% 的关键链路（结算、资产变更、排行榜）

## 适用场景

### ✅ 使用 Stream Broker（JetStream）

1. **结算/资产变更**
   - 战斗结算后发奖
   - 玩家金币/道具变更
   - **要求**：不能丢消息，宕机后可恢复

2. **排行榜更新**
   - 排名变化需要持久化
   - 可回放/补偿
   - **要求**：可靠、有序

3. **邮件/社交事件**
   - 离线期间的邮件/好友请求
   - **要求**：持久化，上线后回放

4. **日志/审计**
   - 需要回放/查询历史
   - **要求**：可审计

### ❌ 不使用 Stream Broker（用 Core）

1. **战斗房间实时通信**（30Hz 状态同步）
   - 要求：低延迟（< 10ms）
   - 可接受：丢中间帧（可快照重建）

2. **玩家指令/操作**
   - 要求：实时响应
   - 可接受：失败重试（业务层处理）

## 当前状态

### ✅ 已完成
- [x] `StreamBroker` 接口定义
- [x] 选项函数（Durable、Queue、MaxAckPending 等）
- [x] 基础实现框架

### 🔧 待完善（需要根据实际 NATS 版本调整）
- [ ] `Enqueue` 方法的具体实现（需要正确的 JetStream API 调用）
- [ ] `Consume` 方法的具体实现（需要正确处理 Ack/Nak/Term）
- [ ] `StreamHandler` 的消息封装（需要能直接调用 Ack/Nak）

## 使用示例（待实现完善后）

```go
// 1. 创建 StreamBroker（需要 NATS 服务器启用 JetStream）
sb, err := nats.NewStreamBroker(nats.WithURLs("nats://127.0.0.1:4222"))
if err != nil {
    log.Fatal(err)
}

// 2. 结算系统：入队结算任务（可靠）
err = sb.Enqueue(ctx, "battle.settle", settleData, 
    broker.WithIdempotencyKey(fmt.Sprintf("battle-%d", battleID)))

// 3. 结算系统：消费队列（可靠处理）
sb.Consume(ctx, "battle.settle", func(ctx context.Context, msg broker.StreamMessage) {
    // 处理结算
    if err := processSettle(msg.Data); err != nil {
        msg.Nak() // 失败，重试
        return
    }
    msg.Ack() // 成功
}, 
    broker.WithDurable("settle-worker"),
    broker.WithStreamQueue("settle-queue"), // 水平扩展
    broker.WithMaxAckPending(100),          // 背压控制
    broker.WithAckWait(30*time.Second))

// 4. 同时可以使用 Core Broker 能力（实时通信）
sb.Publish(ctx, "battle.state", stateData) // 使用 Core 能力
```

## 下一步

1. **根据实际 NATS 版本完善实现**
   - 查看 `go.mod` 中的 `nats.go` 版本
   - 根据对应版本的 JetStream API 完善代码

2. **编写测试用例**
   - 使用嵌入式 NATS Server（已支持 JetStream）
   - 验证 Enqueue/Consume/Ack 流程

3. **实际业务集成**
   - 先在一个结算系统试点
   - 验证可靠性后再扩展到其他场景

## 参考

- [NATS JetStream 文档](https://docs.nats.io/nats-concepts/jetstream)
- [NATS Go JetStream API](https://pkg.go.dev/github.com/nats-io/nats.go/jetstream)
