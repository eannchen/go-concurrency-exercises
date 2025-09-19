package main

import "fmt"

// Processor defines the interface for running a tweet processing example.
// It provides a common contract for different processing strategies.
type Processor interface {
	Process()
}

func main() {
	// A map holds our different processing strategies, keyed by a descriptive name.
	// This makes it easy to add new strategies in the future.
	processors := map[string]Processor{
		"Sequential": &SequentialProcessor{},
		"Concurrent": &ConcurrentProcessor{},
	}

	// Iterate through and run each processor, printing its name first.
	for name, p := range processors {
		fmt.Printf("--- Running %s Processor ---\n", name)
		p.Process()
		fmt.Println()
	}
}
