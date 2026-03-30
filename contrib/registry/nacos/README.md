# Nacos Registry

本包实现 `component/registry` 的 Nacos 适配器。

**特性**
- 支持 group/cluster 维度注册
- 支持权重与默认协议类型
- 支持 Watch 服务变更

**安装**

```bash
go get github.com/byteweap/meta/contrib/registry/nacos
```

**最小用法**

```go
package main

import (
  "context"

  "github.com/byteweap/meta/component/registry"
  nreg "github.com/byteweap/meta/contrib/registry/nacos"
  "github.com/nacos-group/nacos-sdk-go/clients/naming_client"
)

func main() {
  // TODO: 使用 Nacos SDK 初始化 naming client
  var cli naming_client.INamingClient
  reg := nreg.New(cli, nreg.WithGroup("DEFAULT_GROUP"))

  svc := &registry.ServiceInstance{
    ID:        "node-1",
    Name:      "game",
    Version:   "v1",
    Endpoints: []string{"grpc://127.0.0.1:9000"},
  }

  _ = reg.Register(context.Background(), svc)
  _ = reg.Deregister(context.Background(), svc)
}
```

**配置**
- `WithPrefix` 设置前缀路径
- `WithWeight` 设置权重
- `WithCluster` 设置集群名
- `WithGroup` 设置分组
- `WithDefaultKind` 设置默认协议类型

**注意事项**
- `ServiceInstance.Name` 不能为空，否则会返回错误
- Metadata 中可通过 `weight` 覆盖权重
