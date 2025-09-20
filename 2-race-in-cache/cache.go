package main

// KeyStoreCacheLoader is an interface for loading data for the cache.
type KeyStoreCacheLoader interface {
	Load(string) string
}

// Cache defines the interface for a key-value cache.
type Cache interface {
	Get(key string) string
}

// page holds the key-value pair for an entry in the LRU list.
type page struct {
	Key   string
	Value string
}
