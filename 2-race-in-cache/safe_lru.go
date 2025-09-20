package main

import (
	"container/list"
	"sync"
)

// entry is used to signal when a value, currently being loaded, is ready.
type entry struct {
	val   string
	ready chan struct{}
}

// SafeLRUCache is a thread-safe implementation of an LRU cache.
type SafeLRUCache struct {
	cache   map[string]*list.Element
	pages   list.List
	load    func(string) string
	mu      sync.Mutex
	loading map[string]*entry
}

// NewSafeLRUCache creates a new SafeLRUCache.
func NewSafeLRUCache(loader KeyStoreCacheLoader) Cache {
	return &SafeLRUCache{
		load:    loader.Load,
		cache:   make(map[string]*list.Element),
		loading: make(map[string]*entry),
	}
}

// Get is a thread-safe implementation for retrieving from the cache.
func (c *SafeLRUCache) Get(key string) string {
	return c.getOptimized(key)
}

func (c *SafeLRUCache) getOptimized(key string) string {
	c.mu.Lock()
	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		c.mu.Unlock()
		return e.Value.(page).Value
	}

	// Check if their is a goroutine is caching the key
	if req, ok := c.loading[key]; ok {
		// While the key is caching, the goroutine here can just wait for the result
		c.mu.Unlock()
		<-req.ready
		return req.val
	}

	// Mark there is a goroutine will cache the key (it needs to be done before we unlock the `mu`)
	e := &entry{ready: make(chan struct{})}
	c.loading[key] = e

	// Before the long DB operation, we release the `mu`:
	// 1. Other goroutines with the different key can check their cache
	// 2. Other goroutines with the same key will be blocked by `loading` and wait the result
	c.mu.Unlock()
	p := page{key, c.load(key)}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pages.Len() >= CacheSize {
		end := c.pages.Back()
		delete(c.cache, end.Value.(page).Key)
		c.pages.Remove(end)
	}
	c.pages.PushFront(p)
	c.cache[key] = c.pages.Front()

	// After the update, it’s time to let other goroutines with the same key receive the value
	e.val = p.Value
	delete(c.loading, key)
	close(e.ready)

	return p.Value
}

func (c *SafeLRUCache) getWithGlobalLock(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	p := page{key, c.load(key)} // <-- load data from DB (long operation)
	if c.pages.Len() >= CacheSize {
		end := c.pages.Back()
		delete(c.cache, end.Value.(page).Key)
		c.pages.Remove(end)
	}
	c.pages.PushFront(p)
	c.cache[key] = c.pages.Front()

	return p.Value
}
