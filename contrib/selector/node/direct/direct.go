// Package direct 提供直连权重节点
package direct

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/byteweap/wukong/component/selector"
)

const (
	defaultWeight = 100
)

var (
	_ selector.WeightedNode        = (*Node)(nil)
	_ selector.WeightedNodeBuilder = (*Builder)(nil)
)

// Node 表示节点实例
type Node struct {
	selector.Node

	// lastPick 为最近一次选择时间
	lastPick atomic.Int64
}

// Builder 为 direct 节点构建器
type Builder struct{}

// Build 创建节点
func (*Builder) Build(n selector.Node) selector.WeightedNode {
	return &Node{Node: n, lastPick: atomic.Int64{}}
}

// Pick 记录一次选择
func (n *Node) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	n.lastPick.Store(now)
	return func(context.Context, selector.DoneInfo) {}
}

// Weight 返回有效权重
func (n *Node) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return defaultWeight
}

// PickElapsed 返回距上次选择的时间
func (n *Node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - n.lastPick.Load())
}

// Raw 返回原始节点
func (n *Node) Raw() selector.Node {
	return n.Node
}
