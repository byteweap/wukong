// Package wrr 提供加权轮询负载均衡算法
package wrr

import (
	"context"
	"sync"

	"github.com/byteweap/wukong/component/selector"
	defaultselector "github.com/byteweap/wukong/contrib/selector"
	"github.com/byteweap/wukong/contrib/selector/node/direct"
)

const (
	// Name 为 wrr 算法名
	Name = "wrr"
)

var _ selector.Balancer = (*Balancer)(nil)

// Option 为 wrr 构建器选项
type Option func(o *options)

// options 为 wrr 构建器参数
type options struct{}

// Balancer 为 wrr 负载均衡器
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
	lastNodes     []selector.WeightedNode
}

// equalNodes 判断两个节点列表是否相同
func equalNodes(a, b []selector.WeightedNode) bool {
	if len(a) != len(b) {
		return false
	}

	// 构建切片 a 的地址集合
	aMap := make(map[string]bool, len(a))
	for _, node := range a {
		aMap[node.Address()] = true
	}

	// 检查切片 b 的节点是否全部存在
	for _, node := range b {
		if !aMap[node.Address()] {
			return false
		}
	}

	return true
}

// New 创建 wrr 选择器
func New(opts ...Option) selector.Selector {
	return NewBuilder(opts...).Build()
}

// Pick 选择一个节点
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查节点列表是否变化
	if !equalNodes(p.lastNodes, nodes) {
		// 更新 lastNodes
		p.lastNodes = make([]selector.WeightedNode, len(nodes))
		copy(p.lastNodes, nodes)

		// 构建当前节点地址集合用于清理
		currentNodes := make(map[string]bool, len(nodes))
		for _, node := range nodes {
			currentNodes[node.Address()] = true
		}

		// 清理 currentWeight 中的过期节点
		for address := range p.currentWeight {
			if !currentNodes[address] {
				delete(p.currentWeight, address)
			}
		}
	}

	var totalWeight float64
	var selected selector.WeightedNode
	var selectWeight float64

	// nginx wrr 算法参考 http://blog.csdn.net/zhangskd/article/details/50194069
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// 当前权重累加有效权重
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight

	d := selected.Pick()
	return selected, d, nil
}

// NewBuilder 返回带 wrr 的选择器构建器
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

// Builder 为 wrr 构建器
type Builder struct{}

// Build 创建 Balancer
func (b *Builder) Build() selector.Balancer {
	return &Balancer{currentWeight: make(map[string]float64)}
}
