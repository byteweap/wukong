# Meta

<p align="center">
  <img src="assets/logo.svg" alt="Meta logo" width="560">
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/github/go-mod/go-version/byteweap/meta" alt="Go version"></a>
  <a href="https://codecov.io/gh/byteweap/meta"><img src="https://codecov.io/gh/byteweap/meta/branch/master/graph/badge.svg" alt="Codecov"></a>
  <a href="https://goreportcard.com/report/github.com/byteweap/meta"><img src="https://goreportcard.com/badge/github.com/byteweap/meta" alt="Go Report Card"></a>
</p>

Meta 是一个轻量、可插拔、面向长连接游戏/实时业务的 Go 服务框架。它把网关接入、业务逻辑、服务发现、消息通信、玩家位置定位和节点选择拆成清晰的工程边界，让游戏服可以按服务横向扩展，同时保留足够薄的核心运行时。

## 设计目标

- **网关轻量化**：`gate` 只处理 WebSocket 连接、会话管理、消息透传和上下线事件，不承载业务逻辑。
- **业务服务自治**：`mesh` 作为后端逻辑服务运行单元，通过路由注册处理客户端消息、服务间 RPC 和玩家事件。
- **基础设施可替换**：Broker、Registry、Locator、Selector、Logger 都以接口定义在 `component` 下，具体实现放在 `contrib`。
- **面向水平扩展**：通过服务注册发现、节点选择器和玩家位置定位器，把用户流量稳定路由到目标服务节点。
- **保持核心简单**：核心只定义生命周期、通信原语、消息 envelope 和组件边界，限流、分房、匹配、状态同步等策略留给业务层。

## 架构概览

```text
client
  |
  | WebSocket + envelope.IMessage
  v
server/gate
  | 1. 管理 uid -> websocket session
  | 2. 使用 locator 查找玩家所在服务节点
  | 3. 使用 registry + selector 选择可用 mesh 节点
  | 4. 通过 broker 发布消息
  v
component/broker
  |
  | subject: <prefix>.<fromService>.<toService>.<toNodeID>
  v
server/mesh
  | 1. 订阅自身 subject
  | 2. 按 cmd/version 路由业务消息
  | 3. 处理 online/offline/reconnect 系统事件
  | 4. 回包到 gate reply subject
  v
client
```

典型部署中，一个应用进程由 `meta.App` 管理生命周期；进程可以挂载一个或多个 `server.Server` 实例。`gate` 负责客户端入口，`mesh` 负责游戏逻辑或其它后端业务能力。服务实例启动后会注册到 `registry.Registry`，网关通过服务发现监听后端节点变化，并使用 `selector.Selector` 选择目标节点。

## 目录结构

```text
.
├── app.go / option.go          # 应用生命周期、全局选项、服务注册与优雅停止
├── server/                     # 服务运行单元抽象和内置服务
│   ├── server.go               # Server 接口，统一 Gate/Mesh 生命周期
│   ├── gate/                   # WebSocket 网关服务
│   └── mesh/                   # 后端逻辑服务和消息路由
├── component/                  # 核心组件接口
│   ├── broker/                 # pub-sub / request-reply 消息代理抽象
│   ├── registry/               # 服务注册发现抽象
│   ├── locator/                # 玩家位置定位抽象
│   ├── selector/               # 节点选择抽象
│   └── log/                    # 日志抽象
├── contrib/                    # 可选组件实现，按独立模块组织
├── encoding/                   # json/proto/msgpack/yaml/toml/xml 编解码器
├── envelope/                   # 由 idl/envelope 生成的 protobuf 消息
├── internal/cluster/           # 集群 subject、header、系统事件定义
├── pkg/                        # 通用工具包
├── examples/                   # gate、game、gate-client 示例
└── cmd/meta/                   # 命令行脚手架入口
```

## 核心生命周期

`meta.New(...).Run()` 是应用入口。`App` 的主要职责是把服务实例、基础设施组件和进程生命周期串起来：

1. 根据 `ID/Name/Version/Metadata/Endpoint` 构建 `registry.ServiceInstance`。
2. 执行 `PreStart` 钩子。
3. 并发启动所有 `server.Server`。
4. 启动完成后把实例注册到 `Registry`。
5. 执行 `PostStart` 钩子。
6. 监听 `SIGTERM/SIGQUIT/SIGINT` 或外部 `Stop()`。
7. 停止前执行 `PreStop`，从注册中心注销实例。
8. 取消应用上下文，并按 `StopTimeout` 优雅停止所有 server。
9. 执行 `PostStop` 钩子。

`server.Server` 是框架运行单元的最小协议：

```go
type Server interface {
    Kind() Kind
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Endpoint(ctx context.Context) (*url.URL, error)
}
```

这使得 `gate`、`mesh` 或业务自定义 server 都能被同一套生命周期托管。

## Gate 设计

`server/gate` 是客户端接入层，基于 `github.com/olahol/melody` 封装 WebSocket 服务。

启动时，Gate 会：

- 从 `meta.Context` 中读取应用 ID 和服务名。
- 校验必需组件：`Locator`、`Broker`、`Discovery`、`SelectorFunc`。
- 监听 TCP 地址并生成 `ws://host:port` endpoint。
- 注册 Melody 的连接、断开、二进制消息、错误和关闭回调。
- 订阅自己的 broker subject，用来接收 Mesh 回给客户端的消息。

连接建立时，Gate 通过 `UserIDExtractor` 提取 `uid`，默认读取请求参数 `uid`。随后它会：

- 关闭同一用户旧连接，保证单用户单网关会话。
- 在 `Sessions` 中记录 `uid -> *melody.Session`。
- 调用 `locator.Bind(ctx, uid, gateName, gateID)` 绑定用户当前网关节点。
- 向玩家已绑定的业务节点广播 `online` 或 `reconnect` 事件。

客户端二进制消息会被解析为 `envelope.IMessage`：

```proto
message IMessage {
  Header header = 1;
  string service = 2;
  bytes payload = 3;
}
```

Gate 根据 `IMessage.service` 查找目标业务服务：

1. 先通过 `locator.Node(ctx, uid, service)` 查询玩家是否已经绑定到某个业务节点。
2. 如果没有绑定，则通过 `registry.Watch(service)` 维护节点列表，并用 `selector.Selector` 选择一个节点。
3. 使用 `cluster.BuildHeader` 写入 uid、事件类型、reply subject、来源服务和目标服务。
4. 发布到目标 subject：`<prefix>.<gateName>.<service>.<nodeID>`。

Mesh 回包到 Gate 后，Gate 根据 header 中的 uid 找到本地 WebSocket session，并直接 `WriteBinary` 给客户端。

## Mesh 设计

`server/mesh` 是后端逻辑服务。它不直接处理网络连接，而是订阅 broker subject 接收 Gate 或其它服务发来的消息。

Mesh 启动时会校验 `Broker` 和 `Locator`，然后订阅：

```text
<prefix>.*.<appName>.<appID>
```

收到消息后，Mesh 根据 `broker.Message.Reply` 区分两种模式：

- `Reply == ""`：pub-sub 消息，主要用于 Gate 转发的客户端业务消息和系统事件。
- `Reply != ""`：request-reply 消息，用于服务间 RPC。

### 业务路由

客户端业务消息由 `cmd/version` 共同定位处理器。底层路由键使用 `uint64(cmd)<<32 | uint64(version)`，避免字符串拼接开销。

Mesh 支持两类注册方式：

```go
// 推荐：业务函数签名简单，RouteX 内部反射适配
m.RouteX(1, 1, func(ctx *mesh.Context, req *EnterGameReq) {
    ctx.OkResp(&EnterGameResp{})
})

// 热点路径：显式包装，运行期不经过反射调用
m.Route(1, 1, mesh.Wrap(EnterGame))
```

`mesh.Wrap[T]` 会统一完成：

- 从 `envelope.IMessage` 读取 header、service 和 payload。
- 反序列化 protobuf payload 到 `*T`。
- 从对象池获取 `mesh.Context`，减少高频消息分配。
- 调用业务 handler。
- handler 返回后释放 context 到 `sync.Pool`。

`mesh.Context` 提供当前消息的 `UID/Seq/Cmd/Version/Subject/Timestamp/FromService/ToService`，并提供 `OkResp` 和 `ErrResp` 回包方法。若需要在新 goroutine 中使用上下文，应调用 `ctx.Copy()`，因为原始 context 会在 handler 结束后归还对象池。

### 系统事件

Gate 会把用户状态变化广播给玩家已经绑定的业务服务节点。Mesh 内置三个事件钩子：

```go
m.OnlineHandler(func(uid int64) {})
m.OfflineHandler(func(uid int64) {})
m.ReconnectHandler(func(uid int64) {})
```

这些事件使用 broker header 中的 `event` 字段区分，不走业务 `cmd/version` 路由。

### 服务间 RPC

Mesh 还提供 request-reply 能力：

```go
m.RPCRoute("hello", "v1", mesh.WrapRPC(Hello))

func Hello(ctx *mesh.RPCContext, req *HelloReq) ([]byte, string, int) {
    return data, "ok", http.StatusOK
}
```

RPC 路由键为 `cmd.version`。`WrapRPC[T]` 会反序列化请求体并创建 `RPCContext`。返回值约定为：

- `[]byte`：响应数据。
- `string`：提示信息。
- `int`：状态码，`200` 表示成功。

调用方使用：

```go
data, tip, code, err := m.Request(subject, "hello", "v1", payload)
```

## 消息与集群协议

框架把业务消息和集群路由信息分开：

- `envelope.IMessage/OMessage`：客户端与业务层可感知的 protobuf envelope。
- `broker.Header`：集群内部元数据，例如 uid、event、reply、fromService、toService。
- `cluster.Subject`：broker subject 约定。

subject 格式：

```text
<prefix>.<fromService>.<toService>.<toNodeID>
```

例如：

```text
meta.gate.game.game-1
```

这种设计让业务 payload 保持稳定，同时允许集群层独立演进路由、事件和 reply 语义。

## 组件抽象

Meta 的核心组件都以接口形式存在：

- `component/broker.Broker`：提供 `Pub/Sub/Request/Reply/Close`，用于服务间消息通信。
- `component/registry.Registry`：提供 `Register/Deregister/GetService/Watch`，用于服务注册发现。
- `component/locator.Locator`：维护 `uid -> service -> nodeID`，用于玩家粘性路由和状态事件广播。
- `component/selector.Selector`：根据服务实例列表选择目标节点，支持过滤器扩展。
- `component/log.Logger`：统一日志门面，允许替换为 zap/logrus/zerolog 等实现。
- `encoding.Codec`：统一编解码器注册和查找，内置 json、proto、msgpack、yaml、toml、xml。

`contrib` 提供官方可选实现：

- Broker：NATS。
- Locator：Redis。
- Registry：Nacos、Consul、Etcd、Eureka、Polaris、ServiceComb、Zookeeper。
- Selector：Random、RoundRobin、Weighted RoundRobin。
- Logger：Zap、Logrus、Zerolog、Aliyun、Tencent、Fluent。

这些实现通常以独立 Go module 组织，业务可以只引入实际需要的基础设施依赖。

## 快速示例

### 启动 Gate

```go
reg := nacos.New(newNamingClient())
loc := redis.New(goredis.UniversalOptions{
    Addrs: []string{"127.0.0.1:6379"},
}, "meta")
bro, _ := nats.New()

err := meta.New(
    meta.ID("gate-1"),
    meta.Name("gate"),
    meta.Version("v1.0.0"),
    meta.Server(gate.New(
        gate.Addr(":9000"),
        gate.Locator(loc),
        gate.Discovery(reg),
        gate.Broker(bro),
        gate.SelectorFunc(func() selector.Selector {
            return wrr.New()
        }),
    )),
    meta.Registry(reg),
).Run()
```

### 启动 Mesh

```go
type GameServer struct {
    *mesh.Mesh
}

srv := &GameServer{
    Mesh: mesh.New(
        mesh.Broker(bro),
        mesh.Locator(loc),
    ),
}

srv.Route(1, 1, mesh.Wrap(EnterGame))
srv.Route(2, 1, mesh.Wrap(ExitGame))
srv.RPCRoute("hello", "v1", mesh.WrapRPC(Hello))

err := meta.New(
    meta.ID("game-1"),
    meta.Name("game"),
    meta.Version("v1.0.0"),
    meta.Server(srv),
    meta.Registry(reg),
).Run()
```

更多完整代码可以参考：

- `examples/gate`
- `examples/game`
- `examples/gate-client`

## 工程扩展建议

新增业务服务时，推荐继承组合 `*mesh.Mesh`，再把领域状态放在自己的 server 结构体中。例如 `examples/game/internal/server.Server` 把 `roomSpace` 和 `playerSpace` 放在 Mesh 外层，业务 handler 通过自定义 server 访问领域对象。

新增基础设施适配时，建议放在 `contrib/<component>/<provider>` 下，并实现 `component` 中对应接口。核心包不依赖任何具体 provider，这样可以避免把业务项目不需要的 SDK 带进主模块。

高频业务路由建议使用 `Route + mesh.Wrap` 或 `RPCRoute + mesh.WrapRPC`，避免每次请求反射调用；低频或开发期路由可以使用 `RouteX/RPCRouteX` 获得更简洁的注册体验。

如果 handler 内需要异步处理，请复制 `mesh.Context`：

```go
func EnterGame(ctx *mesh.Context, req *EnterGameReq) {
    copied := ctx.Copy()
    go func() {
        _ = copied.UID()
    }()
}
```

## 开发命令

```bash
make test       # 运行所有测试
make tidy       # 整理所有 Go module
make envelope   # 根据 idl/envelope/envelope.proto 重新生成 protobuf
make lint       # 运行 golangci-lint
```

Windows 环境也提供了对应脚本：

```powershell
.\tidy.ps1
.\release.ps1
```

## 当前实现边界

- Gate 文本消息处理仍为预留接口，当前业务链路使用二进制 protobuf envelope。
- Gate 的 request-reply 暂未实现，收到带 reply 的请求会返回 `501 Not Implemented`。
- `mesh.Context.Broadcast` 仍是 TODO。
- Broker 层刻意不内置并发、背压、限流、消息合并等业务策略，这些能力应按具体游戏场景放在 Gate、Mesh 或业务服务中实现。
