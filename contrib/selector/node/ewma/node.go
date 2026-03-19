// Package ewma 提供带统计的权重节点
package ewma

import (
	"context"
	"errors"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/byteweap/wukong/component/selector"
)

const (
	// cost 的平均寿命 在 Tau*ln(2) 处达到半衰期
	tau = int64(time.Millisecond * 600)
	// 如果未收集到统计信息 则对节点增加较大延迟惩罚
	penalty = uint64(time.Microsecond * 100)
)

var (
	_ selector.WeightedNode        = (*Node)(nil)
	_ selector.WeightedNodeBuilder = (*Builder)(nil)
)

// Node 表示节点实例
type Node struct {
	selector.Node

	// 客户端统计数据
	lag       atomic.Int64
	success   atomic.Uint64
	inflight  atomic.Int64
	inflights [200]atomic.Int64
	// 最近一次统计时间
	stamp atomic.Int64
	// 时间窗口内请求次数
	reqs atomic.Int64
	// lastPick 为最近一次选择时间
	lastPick atomic.Int64

	errHandler   func(err error) (isErr bool)
	cachedWeight *atomic.Value
}

type nodeWeight struct {
	value    float64
	updateAt int64
}

// Builder 为 ewma 节点构建器
type Builder struct {
	ErrHandler func(err error) (isErr bool)
}

// Build 创建权重节点
func (b *Builder) Build(n selector.Node) selector.WeightedNode {
	s := &Node{
		Node:         n,
		inflights:    [200]atomic.Int64{},
		errHandler:   b.ErrHandler,
		cachedWeight: &atomic.Value{},
	}
	s.success.Store(1000)
	s.inflight.Store(1)
	return s
}

func (n *Node) health() uint64 {
	return n.success.Load()
}

func (n *Node) load() (load uint64) {
	now := time.Now().UnixNano()
	avgLag := n.lag.Load()
	predict := n.predict(avgLag, now)

	if avgLag == 0 {
		// penalty 为节点刚启动无数据时的惩罚值
		load = penalty * uint64(n.inflight.Load())
		return
	}
	if predict > avgLag {
		avgLag = predict
	}
	// 增加 5ms 以消除不同区域延迟差异
	avgLag += int64(time.Millisecond * 5)
	avgLag = int64(math.Sqrt(float64(avgLag)))
	load = uint64(avgLag) * uint64(n.inflight.Load())
	return load
}

func (n *Node) predict(avgLag int64, now int64) (predict int64) {
	var (
		total    int64
		slowNum  int
		totalNum int
	)
	for i := range n.inflights {
		start := n.inflights[i].Load()
		if start != 0 {
			totalNum++
			lag := now - start
			if lag > avgLag {
				slowNum++
				total += lag
			}
		}
	}
	if slowNum >= (totalNum/2 + 1) {
		predict = total / int64(slowNum)
	}
	return
}

// Pick 选择节点
func (n *Node) Pick() selector.DoneFunc {
	start := time.Now().UnixNano()
	n.lastPick.Store(start)
	n.inflight.Add(1)
	reqs := n.reqs.Add(1)
	slot := reqs % 200
	swapped := n.inflights[slot].CompareAndSwap(0, start)
	return func(_ context.Context, di selector.DoneInfo) {
		if swapped {
			n.inflights[slot].CompareAndSwap(start, 0)
		}
		n.inflight.Add(-1)

		now := time.Now().UnixNano()
		// 获取移动平均系数 w
		stamp := n.stamp.Swap(now)
		td := now - stamp
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(tau))

		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := n.lag.Load()
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		n.lag.Store(lag)

		success := uint64(1000) // 出错时取 1 作为失败值
		if di.Err != nil {
			if n.errHandler != nil && n.errHandler(di.Err) {
				success = 0
			}
			var netErr net.Error
			if errors.Is(context.DeadlineExceeded, di.Err) || errors.Is(context.Canceled, di.Err) ||
				errors.Is(selector.ErrServiceUnavailable, di.Err) || errors.Is(selector.ErrGatewayTimeout, di.Err) || errors.As(di.Err, &netErr) {
				success = 0
			}
		}
		oldSuc := n.success.Load()
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		n.success.Store(success)
	}
}

// Weight 返回有效权重
func (n *Node) Weight() (weight float64) {
	w, ok := n.cachedWeight.Load().(*nodeWeight)
	now := time.Now().UnixNano()
	if !ok || time.Duration(now-w.updateAt) > (time.Millisecond*5) {
		health := n.health()
		load := n.load()
		weight = float64(health*uint64(time.Microsecond)*10) / float64(load)
		n.cachedWeight.Store(&nodeWeight{
			value:    weight,
			updateAt: now,
		})
	} else {
		weight = w.value
	}
	return
}

// PickElapsed 返回距上次选择的时间
func (n *Node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - n.lastPick.Load())
}

// Raw 返回原始节点
func (n *Node) Raw() selector.Node {
	return n.Node
}
