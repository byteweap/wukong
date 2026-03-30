# Selector

## 说明

Selector 用于从节点列表中选择一个可用节点，支持过滤器与权重相关策略。

## 接口概览

核心接口位于 `component/selector`：

- `Select` 选择节点
- `Update` 更新节点列表
- `Nodes` 获取节点列表

辅助接口：

- `Filter` 过滤器函数
- `Node` 节点抽象（`ID/App/Weight/Scheme/Version/Meta`）
- 内置 `Version` 过滤器

## 最小用法

以下示例使用 `contrib/selector/random`：

```go
package main

import (
  "fmt"

  "github.com/byteweap/meta/component/selector"
  "github.com/byteweap/meta/contrib/selector/random"
)

func main() {
  sel := random.NewRandomSelector()
  sel.Update([]selector.Node{
    selector.NewNode("node-a", "", "", "v1", 2, map[string]any{"zone": "sh"}),
    selector.NewNode("node-b", "", "", "v1", 1, map[string]any{"zone": "bj"}),
  })

  n, _ := sel.Select("ignored", selector.Version("v1"))
  fmt.Println(n.ID())
}
```

## 实现（contrib）

- random: `github.com/byteweap/meta/contrib/selector/random`
- roundrobin: `github.com/byteweap/meta/contrib/selector/roundrobin`
- wrr: `github.com/byteweap/meta/contrib/selector/wrr`
