package main

import (
	"container/list"
	"strconv"
	"sync"
	"testing"
	"time"
)

// BenchmarkOptimized tests the optimized implementation with load deduplication
func BenchmarkOptimized(b *testing.B) {
	loader := &Loader{DB: GetMockDB()}
	cache := &SafeLRUCache{
		load:    loader.Load,
		cache:   make(map[string]*list.Element),
		loading: make(map[string]*entry),
	}

	benchmarkCache(b, cache, "optimized")
}

// BenchmarkGlobalLock tests the implementation that holds global lock during load
func BenchmarkGlobalLock(b *testing.B) {
	loader := &Loader{DB: GetMockDB()}
	cache := &SafeLRUCache{
		load:  loader.Load,
		cache: make(map[string]*list.Element),
	}

	benchmarkCache(b, cache, "global-lock")
}

// benchmarkCache runs the actual benchmark with concurrent access
func benchmarkCache(b *testing.B, cache *SafeLRUCache, name string) {
	b.ResetTimer()

	var wg sync.WaitGroup
	start := time.Now()

	// Simulate concurrent access like the test
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i%100) // Use 100 different keys
			if name == "optimized" {
				cache.getOptimized(key)
			} else {
				cache.getWithGlobalLock(key)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	b.ReportMetric(float64(duration.Nanoseconds())/float64(b.N), "ns/op")
}
