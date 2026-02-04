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
