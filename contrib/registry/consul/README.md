# Consul Registry

本包实现 `component/registry` 的 Consul 适配器。

**特性**
- 支持健康检查、心跳与服务标签
- 内置本地缓存与 Watch 机制
- 支持自定义 endpoint 解析与检查配置

**安装**

```bash
go get github.com/byteweap/wukong/contrib/registry/consul
```

**最小用法**

```go
package main

import (
  "context"

  "github.com/byteweap/wukong/component/registry"
  creg "github.com/byteweap/wukong/contrib/registry/consul"
  "github.com/hashicorp/consul/api"
)

func main() {
  cli, _ := api.NewClient(api.DefaultConfig())
  reg := creg.New(cli, creg.WithHealthCheck(true))

  svc := &registry.ServiceInstance{
    ID:        "node-1",
    Name:      "game",
    Version:   "v1",
    Endpoints: []string{"grpc://127.0.0.1:9000"},
    Metadata:  map[string]string{"zone": "sh"},
  }

  _ = reg.Register(context.Background(), svc)
  _ = reg.Deregister(context.Background(), svc)
}
```

**配置**
- `WithHealthCheck` 启用健康检查
- `WithTimeout` 设置获取服务超时
- `WithDatacenter` 设置数据中心
- `WithHeartbeat` 启用心跳
- `WithServiceResolver` 自定义 endpoint 解析
- `WithHealthCheckInterval` 设置检查间隔
- `WithDeregisterCriticalServiceAfter` 设置自动注销时间
- `WithServiceCheck` 追加健康检查
- `WithTags` 设置服务标签

**注意事项**
- `GetService` 优先读取本地缓存，必要时会回源拉取
- `Watch` 会触发后台拉取并广播变更

