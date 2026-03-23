# Zookeeper Registry

本包实现 `component/registry` 的 Zookeeper 适配器。

**特性**
- 使用临时节点注册服务
- 支持 Digest ACL
- 连接恢复后自动重注册

**安装**

```bash
go get github.com/byteweap/wukong/contrib/registry/zookeeper
```

**最小用法**

```go
package main

import (
  "context"
  "time"

  "github.com/byteweap/wukong/component/registry"
  zreg "github.com/byteweap/wukong/contrib/registry/zookeeper"
  "github.com/go-zookeeper/zk"
)

func main() {
  conn, _, _ := zk.Connect([]string{"127.0.0.1:2181"}, 5*time.Second)
  reg := zreg.New(conn, zreg.WithRootPath("/microservices"))

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
- `WithRootPath` 设置注册根路径
- `WithDigestACL` 设置鉴权账号密码

**注意事项**
- 注册节点为临时节点，连接恢复后会自动重建

