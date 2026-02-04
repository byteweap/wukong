package cmap

import (
	"encoding/json"
	"fmt"
	"sync"
)

var SHARD_COUNT = 32

type Stringer interface {
	fmt.Stringer
	comparable
}

// ConcurrentMap 线程安全的分片 map，避免锁竞争
type ConcurrentMap[K comparable, V any] struct {
	shards   []*shardMap[K, V]
	sharding func(key K) uint32
}

// shardMap 分片的数据结构
type shardMap[K comparable, V any] struct {
	items        map[K]V
	sync.RWMutex // 读写锁，保护内部 map
}

func create[K comparable, V any](sharding func(key K) uint32) *ConcurrentMap[K, V] {
	m := &ConcurrentMap[K, V]{
		sharding: sharding,
		shards:   make([]*shardMap[K, V], SHARD_COUNT),
	}
	for i := 0; i < SHARD_COUNT; i++ {
		m.shards[i] = &shardMap[K, V]{items: make(map[K]V)}
	}
	return m
}

// New 创建并发 map
// sharding: 分片函数
func New[K comparable, V any](sharding func(key K) uint32) *ConcurrentMap[K, V] {
	return create[K, V](sharding)
}

// getShard 根据 key 返回对应分片
func (m *ConcurrentMap[K, V]) getShard(key K) *shardMap[K, V] {
	return m.shards[uint(m.sharding(key))%uint(SHARD_COUNT)]
}

// MSet 批量设置键值
func (m *ConcurrentMap[K, V]) MSet(data map[K]V) {
	for key, value := range data {
		shard := m.getShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// Set 设置键值
func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	// 获取分片
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// UpsertCb 回调：锁内调用，禁止访问同 map 的其他键，避免死锁
type UpsertCb[V any] func(exist bool, valueInMap V, newValue V) V

// Upsert 插入或更新：通过回调计算新值
func (m *ConcurrentMap[K, V]) Upsert(key K, value V, cb UpsertCb[V]) (res V) {
	shard := m.getShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

// SetIfAbsent 仅当键不存在时设置
func (m *ConcurrentMap[K, V]) SetIfAbsent(key K, value V) bool {
	// 获取分片
	shard := m.getShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

// Get 获取键值
func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	// 获取分片
	shard := m.getShard(key)
	shard.RLock()
	// 读取分片数据
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Count 返回元素数量
func (m *ConcurrentMap[K, V]) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m.shards[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has 判断键是否存在
func (m *ConcurrentMap[K, V]) Has(key K) bool {
	// 获取分片
	shard := m.getShard(key)
	shard.RLock()
	// 检查是否存在
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove 删除键值
func (m *ConcurrentMap[K, V]) Remove(key K) {
	// 获取分片
	shard := m.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// RemoveCb 锁内回调，返回 true 则删除
type RemoveCb[K any, V any] func(key K, v V, exists bool) bool

// RemoveCb 删除回调
// 锁住分片并执行回调,返回值即回调结果
func (m *ConcurrentMap[K, V]) RemoveCb(key K, cb RemoveCb[K, V]) bool {
	// 获取分片
	shard := m.getShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return remove
}

// Pop 删除并返回键值
func (m *ConcurrentMap[K, V]) Pop(key K) (v V, exists bool) {
	// 获取分片
	shard := m.getShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

// IsEmpty 判断是否为空
func (m *ConcurrentMap[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// Tuple 迭代结果封装
type Tuple[K comparable, V any] struct {
	Key K
	Val V
}

// Iter 返回带缓冲的迭代器
func (m *ConcurrentMap[K, V]) Iter() <-chan Tuple[K, V] {
	cs := snapshot(m)
	total := 0
	for _, c := range cs {
		total += cap(c)
	}
	ch := make(chan Tuple[K, V], total)
	go fanIn(cs, ch)
	return ch
}

// Clear 清空所有键值
func (m *ConcurrentMap[K, V]) Clear() {
	for item := range m.Iter() {
		m.Remove(item.Key)
	}
}

// 创建各分片的快照通道
func snapshot[K comparable, V any](m *ConcurrentMap[K, V]) (cs []chan Tuple[K, V]) {
	// 未初始化时禁止访问
	if len(m.shards) == 0 {
		panic(`cmap.ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	cs = make([]chan Tuple[K, V], SHARD_COUNT)
	wg := sync.WaitGroup{}
	wg.Add(SHARD_COUNT)
	// 遍历分片
	for index, shard := range m.shards {
		go func(index int, shard *shardMap[K, V]) {
			// 遍历键值
			shard.RLock()
			cs[index] = make(chan Tuple[K, V], len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				cs[index] <- Tuple[K, V]{key, val}
			}
			shard.RUnlock()
			close(cs[index])
		}(index, shard)
	}
	wg.Wait()
	return cs
}

// 汇聚多个通道到一个通道
func fanIn[K comparable, V any](cs []chan Tuple[K, V], out chan Tuple[K, V]) {
	wg := sync.WaitGroup{}
	wg.Add(len(cs))
	for _, ch := range cs {
		go func(ch chan Tuple[K, V]) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

// Items 返回所有键值
func (m *ConcurrentMap[K, V]) Items() map[K]V {
	tmp := make(map[K]V)

	// 复制到临时 map
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// IterCb 遍历回调
// 分片内一致，但跨分片不保证一致性
type IterCb[K comparable, V any] func(key K, v V)

// IterCb 回调式遍历，读取成本最低
func (m *ConcurrentMap[K, V]) IterCb(fn IterCb[K, V]) {
	for idx := range m.shards {
		shard := (m.shards)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Keys 返回所有键
func (m *ConcurrentMap[K, V]) Keys() []K {
	count := m.Count()
	ch := make(chan K, count)
	go func() {
		// 遍历分片
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		for _, shard := range m.shards {
			go func(shard *shardMap[K, V]) {
				// 遍历键
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// 生成 key 列表
	keys := make([]K, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// MarshalJSON 序列化
func (m *ConcurrentMap[K, V]) MarshalJSON() ([]byte, error) {
	// 合并分片数据
	tmp := make(map[K]V)

	// 复制到临时 map
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

// UnmarshalJSON 反序列化
func (m *ConcurrentMap[K, V]) UnmarshalJSON(b []byte) (err error) {
	tmp := make(map[K]V)

	// 反序列化为单个 map
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	// 写回并发 map
	for key, val := range tmp {
		m.Set(key, val)
	}
	return nil
}
