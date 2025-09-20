package main

import "container/list"

// CacheSize determines how big the cache can grow.
const CacheSize = 100

// UnsafeLRUCache is the original, race-condition-prone LRU cache.
type UnsafeLRUCache struct {
	cache map[string]*list.Element
	pages list.List
	load  func(string) string
}

// NewUnsafeLRUCache creates a new UnsafeLRUCache.
func NewUnsafeLRUCache(loader KeyStoreCacheLoader) Cache {
	return &UnsafeLRUCache{
		load:  loader.Load,
		cache: make(map[string]*list.Element),
	}
}

// Get gets the key from cache, loads it from the source if needed.
// THIS IMPLEMENTATION CONTAINS A DATA RACE.
func (c *UnsafeLRUCache) Get(key string) string {
	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		return e.Value.(page).Value
	}

	// Miss - load and add to cache
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

// Step 1: The Simultaneous Miss
// Both A and B check the cache. It's empty, so they both get a "cache miss." They proceed to load their respective data from the database.

// Goroutine A now has a page struct: {Key: "Test5", Value: "Test5"}.
// Goroutine B now has a page struct: {Key: "Test8", Value: "Test8"}.

// Step 2: The Unlucky Interleaving (The Race Condition)
// Now, both goroutines try to add their page to the cache.

// Goroutine B runs first: It pushes its page to the list.
// c.pages.PushFront({Key:"Test8", Value:"Test8"})
// The front of the list now contains the data for "Test8".

// Goroutine A runs next: It pushes its page to the list.
// c.pages.PushFront({Key:"Test5", Value:"Test5"})
// The front of the list now contains the data for "Test5". The data for "Test8" is now the second item.

// Goroutine B gets to run again: It updates the map.
// c.cache["Test8"] = c.pages.Front()
// Goroutine B has now incorrectly mapped the key "Test8" to the list element containing the data for "Test5". The cache is now corrupted.
