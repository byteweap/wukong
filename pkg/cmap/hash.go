package cmap

import "fmt"

// FNV1aStr Fowler–Noll–Vo hash 的 1a 变体，用于 fmt.Stringer 入参
func FNV1aStr[K fmt.Stringer](key K) uint32 {
	return FNV1a(key.String())
}

// FNV1a Fowler–Noll–Vo hash 的 1a 变体
func FNV1a(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// Mix64 对 uint64 做高质量混合，并返回低 32 位
// 适合 int64/uint64 类型作为 key 时做分片
func Mix64(key uint64) uint32 {
	key ^= key >> 33
	key *= 0xff51afd7ed558ccd
	key ^= key >> 33
	key *= 0xc4ceb9fe1a85ec53
	key ^= key >> 33
	return uint32(key)
}
