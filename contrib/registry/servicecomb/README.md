# ServiceComb Registry

本包实现 `component/registry` 的 ServiceComb 适配器。

**特性**
- 支持服务注册、发现与 Watch
- 支持心跳上报
- 使用环境变量设置 App 与环境

**安装**

```bash
go get github.com/byteweap/wukong/contrib/registry/servicecomb
```

**最小用法**

```go
package main

import (
  "context"

  "github.com/byteweap/wukong/component/registry"
  sreg "github.com/byteweap/wukong/contrib/registry/servicecomb"
)

func main() {
  // TODO: 使用 ServiceComb 客户端初始化 RegistryClient
  var client sreg.RegistryClient
  reg := sreg.NewRegistry(client)

  svc := &registry.ServiceInstance{
    ID:        "node-1",
    Name:      "game",
    Version:   "v1",
    Endpoints: []string{"rest://127.0.0.1:9000"},
  }

  _ = reg.Register(context.Background(), svc)
  _ = reg.Deregister(context.Background(), svc)
}
```

**注意事项**
- 通过环境变量设置 App 与环境：`CAS_APPLICATION_NAME`、`CAS_ENVIRONMENT_ID`
- `Register` 会处理“服务已存在”的场景并生成实例 ID（若为空）
