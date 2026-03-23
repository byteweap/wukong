package random

import (
	"sort"
	"sync"

	"github.com/byteweap/wukong/component/selector"
	"github.com/byteweap/wukong/pkg/xrand"
)

// Selector 随机选择器（按权重加权随机）
type Selector struct {
	mu          sync.Mutex
	nodes       []selector.Node
	cumulative  []float64
	totalWeight float64
}

var _ selector.Selector = (*Selector)(nil)

const maxRandInt63 = int64(^uint64(0) >> 1)

// New 创建随机选择器实例
func New() *Selector {
	return &Selector{}
}

// Select 选择一个节点（忽略 key）
func (rs *Selector) Select(_ string, filters ...selector.Filter) (selector.Node, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	nodes := rs.nodes
	if len(filters) > 0 {
		nodes = applyFilters(append([]selector.Node(nil), rs.nodes...), filters...)
	}
	if len(nodes) == 0 {
		return nil, selector.ErrNoAvailableNode
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}

	var totalWeight float64
	var cumulative []float64
	if len(filters) == 0 {
		totalWeight = rs.totalWeight
		cumulative = rs.cumulative
	} else {
		cumulative = make([]float64, 0, len(nodes))
		for _, n := range nodes {
			w := n.Weight()
			if w <= 0 {
				w = 1
			}
			totalWeight += w
			cumulative = append(cumulative, totalWeight)
		}
	}
	if totalWeight <= 0 {
		return nil, selector.ErrNoAvailableNode
	}

	randInt := xrand.Int64(0, maxRandInt63)
	r := float64(randInt) / float64(maxRandInt63)
	r *= totalWeight
	idx := sort.Search(len(cumulative), func(i int) bool {
		return cumulative[i] > r
	})
	if idx >= len(nodes) {
		idx = len(nodes) - 1
	}
	return nodes[idx], nil
}

// Update 更新节点列表并重置内部状态
func (rs *Selector) Update(nodes []selector.Node) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.nodes = append([]selector.Node(nil), nodes...)
	rs.cumulative = rs.cumulative[:0]
	rs.totalWeight = 0

	for _, n := range nodes {
		w := n.Weight()
		if w <= 0 {
			w = 1
		}
		rs.totalWeight += w
		rs.cumulative = append(rs.cumulative, rs.totalWeight)
	}
}

// Nodes 返回当前节点列表的副本
func (rs *Selector) Nodes() []selector.Node {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return append([]selector.Node(nil), rs.nodes...)
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
