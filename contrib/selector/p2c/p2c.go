// Package p2c 提供二选一负载均衡算法
package p2c

import (
	"context"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/byteweap/wukong/component/selector"
	defaultselector "github.com/byteweap/wukong/contrib/selector"
	"github.com/byteweap/wukong/contrib/selector/node/ewma"
)

const (
	forcePick = time.Second * 3
	// Name 为 p2c 算法名
	Name = "p2c"
)

var _ selector.Balancer = (*Balancer)(nil)

// Option 为 p2c 构建器选项
type Option func(o *options)

// options 为 p2c 构建器参数
type options struct{}

// New 创建 p2c 选择器
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Balancer 为 p2c 负载均衡器
type Balancer struct {
	mu     sync.Mutex
	r      *rand.Rand
	picked atomic.Bool
}

// prePick 随机选出两个不同节点
func (s *Balancer) prePick(nodes []selector.WeightedNode) (nodeA selector.WeightedNode, nodeB selector.WeightedNode) {
	s.mu.Lock()
	a := s.r.IntN(len(nodes))
	b := s.r.IntN(len(nodes) - 1)
	s.mu.Unlock()
	if b >= a {
		b = b + 1
	}
	nodeA, nodeB = nodes[a], nodes[b]
	return
}

// Pick 选择节点
func (s *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc selector.WeightedNode
	nodeA, nodeB := s.prePick(nodes)
	// meta.Weight 为服务注册时设置的权重
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// 如果低权节点在 forcePick 时间内未被选中过 则强制选择一次
	// 通过强制选择触发成功率与延迟更新
	if upc.PickElapsed() > forcePick && s.picked.CompareAndSwap(false, true) {
		defer s.picked.Store(false)
		pc = upc
	}
	done := pc.Pick()
	return pc, done, nil
}

// NewBuilder 返回带 p2c 的选择器构建器
func NewBuilder(opts ...Option) selector.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}
	return &defaultselector.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &ewma.Builder{},
	}
}

// Builder 为 p2c 构建器
type Builder struct{}

// Build 创建 Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{r: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))}
}
