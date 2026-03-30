# Registry

## 说明

Registry 提供服务注册与发现能力，并支持按服务名监听实例变更。

## 接口概览

核心接口位于 `component/registry`：

- `Register` 注册服务实例
- `Deregister` 注销服务实例
- `GetService` 获取实例列表
- `Watch` 监听实例变更

服务实例 `ServiceInstance` 由以下字段构成：

- `ID` 唯一实例 ID
- `Name` 服务名
- `Version` 版本
- `Metadata` 元数据
- `Endpoints` 访问地址

## 最小用法

```go
package main

import (
  "context"

  "github.com/byteweap/meta/component/registry"
)

func main() {
  var r registry.Registry

  _ = r.Register(context.Background(), &registry.ServiceInstance{
    ID:        "svc-1",
    Name:      "game",
    Version:   "v1.0.0",
    Metadata:  map[string]string{"zone": "sh"},
    Endpoints: []string{"grpc://127.0.0.1:9000"},
  })
}
```

## 实现（contrib）

- etcd: `github.com/byteweap/meta/contrib/registry/etcd`
- consul: `github.com/byteweap/meta/contrib/registry/consul`
- nacos: `github.com/byteweap/meta/contrib/registry/nacos`
- zookeeper: `github.com/byteweap/meta/contrib/registry/zookeeper`
- polaris: `github.com/byteweap/meta/contrib/registry/polaris`
- servicecomb: `github.com/byteweap/meta/contrib/registry/servicecomb`
- eureka: `github.com/byteweap/meta/contrib/registry/eureka`
