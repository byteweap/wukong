package xrand

import (
	"crypto/rand"
	"math/big"
)

// Int [min, max) 随机数生成
// crypto/rand 生成安全的随机数,相比math/rand性能更好，推荐使用
func Int(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(result.Int64()) + min
}

// Int64 [min, max) 随机数生成
// crypto/rand 生成安全的随机数,相比math/rand性能更好，推荐使用
func Int64(min, max int64) int64 {
	if min >= max || max == 0 {
		return max
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return result.Int64() + min
}
