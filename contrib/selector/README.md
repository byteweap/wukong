# Selector Contrib

本模块提供三种常用的节点选择器实现，基于 `github.com/byteweap/meta/component/selector` 的接口规范：

- `random`：按权重加权随机
- `roundrobin`：轮询
- `wrr`：平滑加权轮询（Smooth Weighted Round Robin）

所有实现都 **忽略 key**，仅按节点列表与可选过滤器进行选择。

---

## 安装

```bash
go get github.com/byteweap/meta/contrib/selector
```

分别引用具体实现包即可：

- `github.com/byteweap/meta/contrib/selector/random`
- `github.com/byteweap/meta/contrib/selector/roundrobin`
- `github.com/byteweap/meta/contrib/selector/wrr`

---

## 接口与核心概念

选择器接口位于 `component/selector`：

```go
type Selector interface {
  Select(key string, filters ...Filter) (Node, error)
  Update(nodes []Node)
  Nodes() []Node
}

type Filter func([]Node) []Node

type Node interface {
  ID() string
  App() string
  Weight() float64
  Scheme() string
  Version() string
  Meta() map[string]any
}
```

关键点：

- `Select` 仅与 `nodes` 和 `filters` 相关，当前实现忽略 `key`
- `Update` 会替换内部节点列表（并重置内部状态）
- `Filter` 允许在选择前对节点列表进行裁剪或重排
- `Weight` 用于加权类选择（权重 `<= 0` 会按 `1` 处理）

---

## 最小用法

```go
package main

import (
  "fmt"
  "github.com/byteweap/meta/component/selector"
  "github.com/byteweap/meta/contrib/selector/random"
)

type node struct {
  id     string
  weight float64
  meta   map[string]any
}

func (n node) ID() string           { return n.id }
func (n node) App() string          { return "" }
func (n node) Weight() float64      { return n.weight }
func (n node) Scheme() string       { return "" }
func (n node) Version() string      { return "" }
func (n node) Meta() map[string]any { return n.meta }

func main() {
  sel := random.NewRandomSelector()
  sel.Update([]selector.Node{
    node{id: "node-a", weight: 2, meta: map[string]any{"zone": "sh"}},
    node{id: "node-b", weight: 1, meta: map[string]any{"zone": "bj"}},
  })

  n, err := sel.Select("ignored")
  if err != nil {
    panic(err)
  }
  fmt.Println(n.ID())
}
```

---

## 过滤器

过滤器会在选择前执行，常用于“分区/版本/地区/灰度”等逻辑：

```go
onlySH := func(nodes []selector.Node) []selector.Node {
  out := make([]selector.Node, 0, len(nodes))
  for _, n := range nodes {
    if n.Version() == "v1.0.0"{
      out = append(out, n)
    }
  }
  return out
}

n, err := sel.Select("ignored", onlySH)
```

注意：

- 过滤器执行顺序与传入顺序一致
- 任一过滤器返回空列表会直接导致无可用节点

---

## 实现说明

### random（加权随机）

特性：

- 按权重加权随机选择
- 权重 `<= 0` 会被当作 `1`
- 内部使用 `pkg/xrand` 产生随机数

示例：

```go
import "github.com/byteweap/meta/contrib/selector/random"

sel := random.NewRandomSelector()
```

适用场景：

- 对“均匀打散”要求较高
- 允许统计意义上的公平，但不保证短期平衡

---

### roundrobin（轮询）

特性：

- 严格轮询（无权重）
- 构造时请使用 `NewRoundRobinSelector()`（避免零值 `atomic.Value` 未初始化）

示例：

```go
import "github.com/byteweap/meta/contrib/selector/roundrobin"

sel := roundrobin.NewRoundRobinSelector()
```

适用场景：

- 节点能力相近
- 希望短期内分配尽量均匀

---

### wrr（平滑加权轮询）

特性：

- Smooth Weighted Round Robin
- 权重 `<= 0` 会被当作 `1`
- 节点 `ID()` 应保持唯一（内部索引映射依赖 ID）

示例：

```go
import "github.com/byteweap/meta/contrib/selector/wrr"

sel := wrr.NewWRRSelector()
```

适用场景：

- 节点能力不同
- 需要短期内也相对平滑的权重分配

---

## 线程安全

三个实现均可在并发场景下使用：

- `random` / `wrr`：内部使用互斥锁保护
- `roundrobin`：使用 `atomic.Value` + 原子计数

---

## 常见问题

**1) Select 为什么忽略 key？**  
目前实现只依赖节点列表与过滤器，`key` 预留用于未来一致性哈希等场景。

**2) 权重为 0/负数会怎样？**  
会按 `1` 处理，避免出现“不可选”节点导致总权重为 0。

**3) 过滤器返回空列表会怎样？**  
直接返回 `ErrNoAvailableNode`。

---

## 错误处理

当无可用节点时，会返回 `selector.ErrNoAvailableNode`：

```go
if errors.Is(err, selector.ErrNoAvailableNode) {
  // 处理无可用节点的逻辑
}
```
