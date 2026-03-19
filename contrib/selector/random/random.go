// Package random 提供随机负载均衡算法
package random

import (
	"context"
	"math/rand/v2"

	"github.com/byteweap/wukong/component/selector"
	defaultselector "github.com/byteweap/wukong/contrib/selector"
	"github.com/byteweap/wukong/contrib/selector/node/direct"
)

const (
	// Name 为 random 算法名
	Name = "random"
)

var _ selector.Balancer = (*Balancer)(nil)

// Option 为 random 构建器选项
type Option func(o *options)

// options 为 random 构建器参数
type options struct{}

// Balancer 为随机负载均衡器
type Balancer struct{}

// New 创建随机选择器
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick 选择一个节点
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	cur := rand.IntN(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder 返回带 random 的选择器构建器
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &defaultselector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &direct.Builder{},
	}
}

// Builder 为 random 构建器
type Builder struct{}

// Build 创建 Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{}
}
