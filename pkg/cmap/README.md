# cmap

`cmap` 是一个并发安全的分片 map，实现方式为“多个分片 + 每分片一把读写锁”。它适用于高并发访问场景，可有效减少单锁竞争。

## 快速使用

`cmap` 仅保留一个 `New`，调用时需要显式提供分片（哈希）函数：

```go
// string 作为 key
m := cmap.New[string, *Session](cmap.FNV1a)

// 写入
m.Set(playerID, session)

// 读取
sess, ok := m.Get(playerID)
if ok {
    // ...
}

// 删除
m.Remove(playerID)
```

对于实现了 `fmt.Stringer` 的类型：

```go
// K 实现 fmt.Stringer
m := cmap.New[MyKey, *Session](cmap.FNV1aStr[MyKey])
```

## int64 用户 ID 的推荐分片

如果 key 是 `int64` 类型的用户 ID，建议先做 64 位混合（mix）再取低 32 位，以改善分布：

```go
func mix64(x uint64) uint64 {
    x ^= x >> 33
    x *= 0xff51afd7ed558ccd
    x ^= x >> 33
    x *= 0xc4ceb9fe1a85ec53
    x ^= x >> 33
    return x
}

func shardingInt64(id int64) uint32 {
    return uint32(mix64(uint64(id)))
}

m := cmap.New[int64, *Session](shardingInt64)
```

## 内置分片算法（hash.go）

- `FNV1a(key string) uint32`：FNV-1a 32 位哈希
- `FNV1aStr[K fmt.Stringer](key K) uint32`：对 `Stringer` 的 `String()` 结果做 FNV-1a

## 常用 API

- `Set(key K, value V)`
- `Get(key K) (V, bool)`
- `SetIfAbsent(key K, value V) bool`
- `Upsert(key K, value V, cb UpsertCb[V]) V`
- `Remove(key K)`
- `RemoveCb(key K, cb RemoveCb[K, V]) bool`
- `Pop(key K) (V, bool)`
- `Has(key K) bool`
- `Count() int`
- `IsEmpty() bool`
- `Keys() []K`
- `Items() map[K]V`
- `Iter() <-chan Tuple[K, V)`
- `IterCb(fn IterCb[K, V])`
- `Clear()`
- `MarshalJSON() ([]byte, error)`
- `UnmarshalJSON(b []byte) error`

## 性能提醒

`sharding` 算法会直接影响分片分布与锁竞争，进而显著影响 `cmap` 性能。请务必结合实际 key 分布与业务并发特征做好基准测试后再选型。

## 适用场景对比：cmap vs sync.Map

### cmap 更适合
- 写入并发高、读写比例接近、写多读少 的场景
- key 分布相对均匀，热点不明显
- 需要更可控的锁竞争模型

### sync.Map 更适合
- 读多写少、key 生命周期较长的场景
- 热点 key 被频繁读取
- 业务以缓存访问为主，写入较少

## 结合长链接网关场景的建议

对于“玩家 ID -> Session”的网关管理场景：
- 玩家上线/下线是低频写入
- 每条消息都需要查找 Session（高频读取）
- Session 生命周期较长（1–30 分钟）

因此整体更偏“读多写少 + key 生命周期长”的特征，通常 `sync.Map` 的读路径更有优势。

如果你需要更强的控制（例如统计/全量遍历/一致性快照），或者实际写入比例更高，再考虑使用 `cmap`。
