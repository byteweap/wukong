// Package defaultselector 提供默认选择器实现
package defaultselector

import (
	"context"
	"sync/atomic"

	"github.com/byteweap/wukong/component/selector"
)

var (
	_ selector.Rebalancer = (*Default)(nil)
	_ selector.Builder    = (*DefaultBuilder)(nil)
)

// Default 是组合选择器
type Default struct {
	NodeBuilder selector.WeightedNodeBuilder
	Balancer    selector.Balancer

	nodes atomic.Value
}

// Select 选择一个节点
func (d *Default) Select(ctx context.Context, opts ...selector.SelectOption) (selected selector.Node, done selector.DoneFunc, err error) {
	var (
		options    selector.SelectOptions
		candidates []selector.WeightedNode
	)
	nodes, ok := d.nodes.Load().([]selector.WeightedNode)
	if !ok {
		return nil, nil, selector.ErrNoAvailable
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NodeFilters) > 0 {
		newNodes := make([]selector.Node, len(nodes))
		for i, wc := range nodes {
			newNodes[i] = wc
		}
		for _, filter := range options.NodeFilters {
			newNodes = filter(ctx, newNodes)
		}
		candidates = make([]selector.WeightedNode, len(newNodes))
		for i, n := range newNodes {
			candidates[i] = n.(selector.WeightedNode)
		}
	} else {
		candidates = nodes
	}

	if len(candidates) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	wn, done, err := d.Balancer.Pick(ctx, candidates)
	if err != nil {
		return nil, nil, err
	}
	p, ok := selector.FromPeerContext(ctx)
	if ok {
		p.Node = wn.Raw()
	}
	return wn.Raw(), done, nil
}

// Apply 更新节点信息
func (d *Default) Apply(nodes []selector.Node) {
	weightedNodes := make([]selector.WeightedNode, 0, len(nodes))
	for _, n := range nodes {
		weightedNodes = append(weightedNodes, d.NodeBuilder.Build(n))
	}
	// TODO 不删除未变化节点
	d.nodes.Store(weightedNodes)
}

// DefaultBuilder 是默认构建器
type DefaultBuilder struct {
	Node     selector.WeightedNodeBuilder
	Balancer selector.BalancerBuilder
}

// Build 创建 Selector
func (db *DefaultBuilder) Build() selector.Selector {
	return &Default{
		NodeBuilder: db.Node,
		Balancer:    db.Balancer.Build(),
	}
}
