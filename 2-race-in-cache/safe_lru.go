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
func (c *SafeLRUCache) Get(key string) string {
	return c.getOptimized(key)
	// return c.getWithGlobalLock(key)
}

func (c *SafeLRUCache) getOptimized(key string) string {
	// TODO: Implement optimized get
	return c.getWithGlobalLock(key)
}

func (c *SafeLRUCache) getWithGlobalLock(key string) string {
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
