package cmap

import (
	"strconv"
	"testing"
)

func TestFNV1a(t *testing.T) {

	countMap := make(map[uint32]int)

	var num uint64
	for range 10000 {
		num += 5
		val := FNV1a(strconv.FormatUint(num, 10)) % 32
		countMap[val]++
	}
	for k, v := range countMap {
		t.Logf("%v: %v", k, v)
	}
}

// TestMix64 测试Mix64函数hash后是否均匀分布
func TestMix64(t *testing.T) {

	countMap := make(map[uint32]int)

	var num uint64
	for range 10000 {
		num += 5
		val := Mix64(num) % 32
		countMap[val]++
	}
	for k, v := range countMap {
		t.Logf("%v: %v", k, v)
	}
}
