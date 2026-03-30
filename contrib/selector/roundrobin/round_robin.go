package roundrobin

import (
	"sync/atomic"

	"github.com/byteweap/meta/component/selector"
)

// Selector 轮询选择器
type Selector struct {
	nodes atomic.Value // []selector.Node
	next  uint64
}

var _ selector.Selector = (*Selector)(nil)

// New 创建轮询选择器实例
func New() *Selector {
	rr := &Selector{}
	rr.nodes.Store([]selector.Node{})
	return rr
}

// Select 选择一个节点（忽略 key）
func (rr *Selector) Select(_ string, filters ...selector.Filter) (selector.Node, error) {
	v := rr.nodes.Load()
	nodes := v.([]selector.Node)
	if len(filters) > 0 {
		nodes = applyFilters(append([]selector.Node(nil), nodes...), filters...)
	}
	if len(nodes) == 0 {
		return nil, selector.ErrNoAvailableNode
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	idx := atomic.AddUint64(&rr.next, 1)
	return nodes[(idx-1)%uint64(len(nodes))], nil
}

// Update 更新节点列表并重置轮询位置
func (rr *Selector) Update(nodes []selector.Node) {
	rr.nodes.Store(append([]selector.Node(nil), nodes...))
	atomic.StoreUint64(&rr.next, 0)
}

// Nodes 返回当前节点列表的副本
func (rr *Selector) Nodes() []selector.Node {
	v := rr.nodes.Load()
	nodes := v.([]selector.Node)
	return append([]selector.Node(nil), nodes...)
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
