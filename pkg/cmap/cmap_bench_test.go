package cmap

import (
	"strconv"
	"sync"
	"testing"
)

type Integer int

func (i Integer) String() string {
	return strconv.Itoa(int(i))
}

func BenchmarkItems(b *testing.B) {
	m := New[string, Animal](FNV1a)

	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkItemsInteger(b *testing.B) {
	m := New[Integer, Animal](FNV1aStr)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set((Integer)(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}
func directSharding(key uint32) uint32 {
	return key
}

func BenchmarkItemsInt(b *testing.B) {
	m := New[uint32, Animal](directSharding)

	for i := 0; i < 10000; i++ {
		m.Set((uint32)(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Items()
	}
}

func BenchmarkMarshalJson(b *testing.B) {
	m := New[string, Animal](FNV1a)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		_, err := m.MarshalJSON()
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}

func BenchmarkHashMix64(b *testing.B) {
	var key uint64 = 123456789
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key += uint64(i)
		_ = Mix64(key)
	}
}

func BenchmarkSingleInsertAbsent(b *testing.B) {
	b.Run("cmap", func(b *testing.B) {
		m := New[string, string](FNV1a)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Set(strconv.Itoa(i), "value")
		}
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Store(strconv.Itoa(i), "value")
		}
	})
}

func BenchmarkSingleInsertPresent(b *testing.B) {
	b.Run("cmap", func(b *testing.B) {
		m := New[string, string](FNV1a)
		m.Set("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Set("key", "value")
		}
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		m.Store("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m.Store("key", "value")
		}
	})
}

func benchmarkMultiInsertDifferent(b *testing.B) {
	m := New[string, string](FNV1a)
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferentSyncMap(b *testing.B) {
	var m sync.Map
	finished := make(chan struct{}, b.N)
	_, set := GetSetSyncMap[string, string](&m, finished)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent(b *testing.B) {
	b.Run("shards_1", func(b *testing.B) {
		runWithShards(benchmarkMultiInsertDifferent, b, 1)
	})
	b.Run("shards_16", func(b *testing.B) {
		runWithShards(benchmarkMultiInsertDifferent, b, 16)
	})
	b.Run("shards_32", func(b *testing.B) {
		runWithShards(benchmarkMultiInsertDifferent, b, 32)
	})
	b.Run("shards_64", func(b *testing.B) {
		runWithShards(benchmarkMultiInsertDifferent, b, 64)
	})
	b.Run("shards_256", func(b *testing.B) {
		runWithShards(benchmarkMultiInsertDifferent, b, 256)
	})
}

func BenchmarkMultiInsertSame(b *testing.B) {
	b.Run("cmap", func(b *testing.B) {
		m := New[string, string](FNV1a)
		finished := make(chan struct{}, b.N)
		_, set := GetSet(m, finished)
		m.Set("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go set("key", "value")
		}
		for i := 0; i < b.N; i++ {
			<-finished
		}
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		finished := make(chan struct{}, b.N)
		_, set := GetSetSyncMap[string, string](&m, finished)
		m.Store("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go set("key", "value")
		}
		for i := 0; i < b.N; i++ {
			<-finished
		}
	})
}

func BenchmarkMultiGetSame(b *testing.B) {
	b.Run("cmap", func(b *testing.B) {
		m := New[string, string](FNV1a)
		finished := make(chan struct{}, b.N)
		get, _ := GetSet(m, finished)
		m.Set("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go get("key", "value")
		}
		for i := 0; i < b.N; i++ {
			<-finished
		}
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		finished := make(chan struct{}, b.N)
		get, _ := GetSetSyncMap[string, string](&m, finished)
		m.Store("key", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go get("key", "value")
		}
		for i := 0; i < b.N; i++ {
			<-finished
		}
	})
}

func benchmarkMultiGetSetDifferent(b *testing.B) {
	m := New[string, string](FNV1a)
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i-1), "value")
		go get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent(b *testing.B) {
	b.Run("cmap_shards_1", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetDifferent, b, 1)
	})
	b.Run("cmap_shards_16", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetDifferent, b, 16)
	})
	b.Run("cmap_shards_32", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetDifferent, b, 32)
	})
	b.Run("cmap_shards_256", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetDifferent, b, 256)
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		finished := make(chan struct{}, 2*b.N)
		get, set := GetSetSyncMap[string, string](&m, finished)
		m.Store("-1", "value")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go set(strconv.Itoa(i-1), "value")
			go get(strconv.Itoa(i), "value")
		}
		for i := 0; i < 2*b.N; i++ {
			<-finished
		}
	})
}

func benchmarkMultiGetSetBlock(b *testing.B) {
	m := New[string, string](FNV1a)
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go set(strconv.Itoa(i%100), "value")
		go get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock(b *testing.B) {
	b.Run("cmap_shards_1", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetBlock, b, 1)
	})
	b.Run("cmap_shards_16", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetBlock, b, 16)
	})
	b.Run("cmap_shards_32", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetBlock, b, 32)
	})
	b.Run("cmap_shards_256", func(b *testing.B) {
		runWithShards(benchmarkMultiGetSetBlock, b, 256)
	})

	b.Run("syncmap", func(b *testing.B) {
		var m sync.Map
		finished := make(chan struct{}, 2*b.N)
		get, set := GetSetSyncMap[string, string](&m, finished)
		for i := 0; i < b.N; i++ {
			m.Store(strconv.Itoa(i%100), "value")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go set(strconv.Itoa(i%100), "value")
			go get(strconv.Itoa(i%100), "value")
		}
		for i := 0; i < 2*b.N; i++ {
			<-finished
		}
	})
}

func GetSet[K comparable, V any](m *ConcurrentMap[K, V], finished chan struct{}) (set func(key K, value V), get func(key K, value V)) {
	return func(key K, value V) {
			for i := 0; i < 10; i++ {
				m.Get(key)
			}
			finished <- struct{}{}
		}, func(key K, value V) {
			for i := 0; i < 10; i++ {
				m.Set(key, value)
			}
			finished <- struct{}{}
		}
}

func GetSetSyncMap[K comparable, V any](m *sync.Map, finished chan struct{}) (get func(key K, value V), set func(key K, value V)) {
	get = func(key K, value V) {
		for i := 0; i < 10; i++ {
			m.Load(key)
		}
		finished <- struct{}{}
	}
	set = func(key K, value V) {
		for i := 0; i < 10; i++ {
			m.Store(key, value)
		}
		finished <- struct{}{}
	}
	return
}

func runWithShards(bench func(b *testing.B), b *testing.B, shardsCount int) {
	oldShardsCount := SHARD_COUNT
	SHARD_COUNT = shardsCount
	bench(b)
	SHARD_COUNT = oldShardsCount
}

func BenchmarkKeys(b *testing.B) {
	m := New[string, Animal](FNV1a)

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		m.Set(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for i := 0; i < b.N; i++ {
		m.Keys()
	}
}
