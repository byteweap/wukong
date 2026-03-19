package selector

import (
	"context"
	"time"
)

// Balancer 定义负载均衡器接口
type Balancer interface {
	Pick(ctx context.Context, nodes []WeightedNode) (selected WeightedNode, done DoneFunc, err error)
}

// BalancerBuilder 构建 Balancer
type BalancerBuilder interface {
	Build() Balancer
}

// WeightedNode 提供运行时权重
type WeightedNode interface {
	Node

	// Raw 返回原始节点
	Raw() Node

	// Weight 返回运行时权重
	Weight() float64

	// Pick 记录一次选择
	Pick() DoneFunc

	// PickElapsed 返回距上次选择的时间
	PickElapsed() time.Duration
}

// WeightedNodeBuilder 构建 WeightedNode
type WeightedNodeBuilder interface {
	Build(Node) WeightedNode
}
