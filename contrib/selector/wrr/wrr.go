package wrr

import (
	"sync"

	"github.com/byteweap/wukong/component/selector"
)

// Selector Weighted Round Robin 加权轮询选择器
// 注意：选择不依赖 key，仅按权重轮询
type Selector struct {
	mu          sync.Mutex
	nodes       []selector.Node
	state       []state
	totalWeight float64
	index       map[string]int
}

type state struct {
	node    selector.Node
	weight  float64
	current float64
}

var _ selector.Selector = (*Selector)(nil)

// New 创建 WRRSelector 实例
func New() *Selector {
	return &Selector{}
}

// Select 选择一个节点（忽略 key）
func (ws *Selector) Select(_ string, filters ...selector.Filter) (selector.Node, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	nodes := ws.nodes
	if len(filters) > 0 {
		nodes = applyFilters(append([]selector.Node(nil), ws.nodes...), filters...)
	}
	if len(nodes) == 0 {
		return nil, selector.ErrNoAvailableNode
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}

	if len(filters) == 0 {
		best := 0
		for i := range ws.state {
			ws.state[i].current += ws.state[i].weight
			if i == 0 || ws.state[i].current > ws.state[best].current {
				best = i
			}
		}
		ws.state[best].current -= ws.totalWeight
		return ws.state[best].node, nil
	}

	candidates := make([]int, 0, len(nodes))
	var totalWeight float64
	for _, n := range nodes {
		idx, ok := ws.index[n.ID()]
		if !ok {
			continue
		}
		candidates = append(candidates, idx)
		totalWeight += ws.state[idx].weight
	}
	if len(candidates) == 0 || totalWeight <= 0 {
		return nil, selector.ErrNoAvailableNode
	}

	best := candidates[0]
	for _, idx := range candidates {
		ws.state[idx].current += ws.state[idx].weight
		if ws.state[idx].current > ws.state[best].current {
			best = idx
		}
	}
	ws.state[best].current -= totalWeight
	return ws.state[best].node, nil
}

// Update 更新节点列表并重置内部状态
func (ws *Selector) Update(nodes []selector.Node) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.nodes = append([]selector.Node(nil), nodes...)
	ws.state = make([]state, 0, len(nodes))
	ws.totalWeight = 0
	ws.index = make(map[string]int, len(nodes))

	for i, n := range nodes {
		w := n.Weight()
		if w <= 0 {
			w = 1
		}
		ws.state = append(ws.state, state{
			node:   n,
			weight: w,
		})
		ws.index[n.ID()] = i
		ws.totalWeight += w
	}
}

// Nodes 返回当前节点列表的副本
func (ws *Selector) Nodes() []selector.Node {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return append([]selector.Node(nil), ws.nodes...)
}

func applyFilters(nodes []selector.Node, filters ...selector.Filter) []selector.Node {
	out := nodes
	for _, f := range filters {
		if f == nil {
			continue
		}
		out = f(out)
		if len(out) == 0 {
			return out
		}
	}
	return out
}
