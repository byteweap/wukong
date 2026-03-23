# Etcd Registry

本包实现 `component/registry` 的 etcd 适配器。

**特性**
- 基于租约 + KeepAlive 实现注册心跳
- 支持自定义命名空间与重试策略
- 支持 Watch 服务变更

**安装**

```bash
go get github.com/byteweap/wukong/contrib/registry/etcd
```

**最小用法**

```go
package main

import (
  "context"
  "time"

  "github.com/byteweap/wukong/component/registry"
  ereg "github.com/byteweap/wukong/contrib/registry/etcd"
  clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
  cli, _ := clientv3.New(clientv3.Config{
    Endpoints: []string{"127.0.0.1:2379"},
  })
  reg := ereg.New(cli, ereg.RegisterTTL(15*time.Second))

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
- `Context` 设置注册中心上下文
- `Namespace` 设置命名空间，默认 `/microservices`
- `RegisterTTL` 设置租约 TTL
- `MaxRetry` 设置心跳重试次数

**注意事项**
- 注册 key 结构：`<namespace>/<service>/<id>`
- 注销会停止对应心跳
