package main

import (
	"container/list"
	"flag"
	"fmt"
	"time"
)

var (
	useGlobalLock = flag.Bool("global-lock", false, "Use global lock implementation")
)

func main() {
	flag.Parse()

	var cache Cache
	loader := &Loader{DB: GetMockDB()}

	if *useGlobalLock {
		cache = &globalLockWrapper{&SafeLRUCache{
			load:  loader.Load,
			cache: make(map[string]*list.Element),
		}}
		fmt.Println("Using Global Lock Implementation")
	} else {
		cache = &optimizedWrapper{&SafeLRUCache{
			load:  loader.Load,
			cache: make(map[string]*list.Element),
		}}
		fmt.Println("Using Optimized Implementation")
	}

	fmt.Println("Running server simulation...")
	start := time.Now()
	RunMockServer(cache, nil)
	duration := time.Since(start)

	fmt.Printf("Simulation finished in %s.\n", duration)
	fmt.Printf("DB calls: %d\n", loader.DB.Calls)
	fmt.Printf("Average time per call: %v\n", duration/time.Duration(cycles*callsPerCycle))
}

// optimizedWrapper wraps SafeLRUCache to use optimized method
type optimizedWrapper struct {
	*SafeLRUCache
}

func (w *optimizedWrapper) Get(key string) string {
	return w.getOptimized(key)
}

// globalLockWrapper wraps SafeLRUCache to use global lock method
type globalLockWrapper struct {
	*SafeLRUCache
}

func (w *globalLockWrapper) Get(key string) string {
	return w.getWithGlobalLock(key)
}
