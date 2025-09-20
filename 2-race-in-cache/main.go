package main

import (
	"fmt"
	"time"
)

func main() {
	loader := &Loader{DB: GetMockDB()}

	// You can switch this to NewUnsafeLRUCache to see the race condition.
	cache := NewSafeLRUCache(loader)

	fmt.Println("Running server simulation...")
	start := time.Now()
	RunMockServer(cache, nil)
	fmt.Printf("Simulation finished in %s.\n", time.Since(start))
}
