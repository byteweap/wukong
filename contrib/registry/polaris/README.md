# Polaris Registry

本包实现 `component/registry` 的 Polaris 适配器。

**特性**
- 支持配置命名空间、Token、权重、TTL 等
- 支持心跳上报与实例隔离
- 支持 Watch 服务变更

**安装**

```bash
go get github.com/byteweap/wukong/contrib/registry/polaris
```

**最小用法**

```go
package main

import (
  "context"

  "github.com/byteweap/wukong/component/registry"
  preg "github.com/byteweap/wukong/contrib/registry/polaris"
  "github.com/polarismesh/polaris-go/pkg/config"
)

func main() {
  // TODO: 使用 polaris-go 构造配置
  var conf config.Configuration
  reg := preg.NewRegistryWithConfig(conf, preg.WithNamespace("default"))

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
- `WithNamespace` 设置命名空间
- `WithServiceToken` 设置服务访问 Token
- `WithProtocol` 设置协议类型
- `WithWeight` 设置权重
- `WithHealthy` 设置健康状态
- `WithIsolate` 设置隔离
- `WithTTL` 设置 TTL
- `WithTimeout` 设置超时时间
- `WithRetryCount` 设置重试次数
- `WithHeartbeat` 设置是否启用心跳

**注意事项**
- 注册时服务名会拼接协议：`<Name><scheme>`
- 启用心跳时 `TTL` 必须设置
