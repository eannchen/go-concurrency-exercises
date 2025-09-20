package main

import (
	"strconv"
	"testing"
)

// TestSafeCache_Correctness tests the safe implementation for correctness and races.
func TestSafeCache_Correctness(t *testing.T) {
	loader := &Loader{DB: GetMockDB()}
	cache := NewSafeLRUCache(loader)

	RunMockServer(cache, t)

	if loader.DB.Calls > callsPerCycle {
		t.Errorf("Too many DB calls: %v, expected around %v", loader.DB.Calls, callsPerCycle)
	}
}

// TestSafeCache_LRU_Logic verifies the LRU eviction logic of the safe cache.
func TestSafeCache_LRU_Logic(t *testing.T) {
	loader := &Loader{DB: GetMockDB()}
	cache := NewSafeLRUCache(loader)

	// Sequentially fill the cache to have a predictable LRU order.
	for i := 0; i < 100; i++ {
		cache.Get("key" + strconv.Itoa(i))
	}

	// At this point, key0 is the least recently used.
	// Access key0 to make it the MOST recently used.
	cache.Get("key0")

	// Now, key1 is the least recently used.
	// Add one more item, which should evict key1.
	cache.Get("key100")

	// To check the internal state, we need a type assertion.
	safeCache := cache.(*SafeLRUCache)
	safeCache.mu.Lock()
	defer safeCache.mu.Unlock()

	// Check that key1 was evicted.
	if _, ok := safeCache.cache["key1"]; ok {
		t.Errorf("key1 should have been evicted, but it was not")
	}

	// Check that key0 was NOT evicted.
	if _, ok := safeCache.cache["key0"]; !ok {
		t.Errorf("key0 was evicted, but it should not have been")
	}
}
