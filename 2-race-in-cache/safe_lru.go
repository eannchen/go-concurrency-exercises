package main

import (
	"container/list"
	"sync"
)

// SafeLRUCache is a thread-safe implementation of an LRU cache.
type SafeLRUCache struct {
	cache    map[string]*list.Element
	pages    list.List
	load     func(string) string
	mu       sync.Mutex
	loadLock sync.Mutex
}

// NewSafeLRUCache creates a new SafeLRUCache.
func NewSafeLRUCache(loader KeyStoreCacheLoader) Cache {
	return &SafeLRUCache{
		load:  loader.Load,
		cache: make(map[string]*list.Element),
	}
}

// Get is a thread-safe implementation for retrieving from the cache.
// It is optimized to not hold the lock during slow database calls.
func (c *SafeLRUCache) Get(key string) string {
	return c.get(key)
	// return c.getWithBadPerformance(key)
}

func (c *SafeLRUCache) get(key string) string {
	c.mu.Lock()
	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		c.mu.Unlock()
		return e.Value.(page).Value
	}

	// Release the lock, so other goroutines can still read from the cache.
	c.mu.Unlock()
	value := c.load(key) // <-- SLOW OPERATION

	// Acquire the lock again because we are about to modify the cache
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check the cache again.
	// It's possible another goroutine fetched and stored the same key while we were busy loading it from the database
	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	p := page{key, value}
	if c.pages.Len() >= CacheSize {
		end := c.pages.Back()
		delete(c.cache, end.Value.(page).Key)
		c.pages.Remove(end)
	}
	c.pages.PushFront(p)
	c.cache[key] = c.pages.Front()

	return p.Value
}

func (c *SafeLRUCache) getWithBadPerformance(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	p := page{key, c.load(key)}
	if c.pages.Len() >= CacheSize {
		end := c.pages.Back()
		delete(c.cache, end.Value.(page).Key)
		c.pages.Remove(end)
	}
	c.pages.PushFront(p)
	c.cache[key] = c.pages.Front()

	return p.Value
}
