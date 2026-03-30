# Eureka Registry

本包实现 `component/registry` 的 Eureka 适配器。

**特性**
- 支持心跳与服务列表刷新
- endpoints 映射为 Eureka Instance
- 支持 Watch 服务变更

**安装**

```bash
go get github.com/byteweap/meta/contrib/registry/eureka
```

**最小用法**

```go
package main

import (
  "context"

  "github.com/byteweap/meta/component/registry"
  ereg "github.com/byteweap/meta/contrib/registry/eureka"
)

func main() {
  reg, _ := ereg.New([]string{"http://127.0.0.1:8761"})

  svc := &registry.ServiceInstance{
    ID:        "node-1",
    Name:      "game",
    Version:   "v1",
    Endpoints: []string{"http://127.0.0.1:9000"},
  }

  _ = reg.Register(context.Background(), svc)
  _ = reg.Deregister(context.Background(), svc)
}
```

**配置**
- `WithContext` 设置注册中心上下文
- `WithHeartbeat` 设置心跳间隔
- `WithRefresh` 设置刷新间隔
- `WithEurekaPath` 设置 Eureka 接口路径

**注意事项**
- `Endpoints` 将被解析为 Eureka Instance，Metadata 可覆盖 `securePort`、`homePageURL`、`statusPageURL`、`healthCheckURL`

