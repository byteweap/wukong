package xrand

import (
	"crypto/rand"
	"math/big"
)

// Int [minValue, maxValue) 随机数生成
// crypto/rand 生成安全的随机数,相比math/rand性能更好，推荐使用
func Int(minValue, maxValue int) int {
	if minValue >= maxValue || maxValue == 0 {
		return maxValue
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(maxValue-minValue)))
	return int(result.Int64()) + minValue
}

// Int64 [minValue, maxValue) 随机数生成
// crypto/rand 生成安全的随机数,相比math/rand性能更好，推荐使用
func Int64(minValue, maxValue int64) int64 {
	if minValue >= maxValue || maxValue == 0 {
		return maxValue
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(maxValue-minValue))
	return result.Int64() + minValue
}
