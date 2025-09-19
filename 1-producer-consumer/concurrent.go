package main

import (
	"fmt"
	"sync"
	"time"
)

// ConcurrentProcessor handles processing tweets concurrently using goroutines.
type ConcurrentProcessor struct{}

const numConsumers = 2 // Defines how many concurrent consumers we want (fan-out)

// Process runs the concurrent pipeline.
func (n *ConcurrentProcessor) Process() {
	start := time.Now()
	stream := GetMockStream()

	// This is a good practice to prevent memory leak,
	// even though we know consumer will receive all values.
	done := make(chan struct{})
	defer close(done)

	tweets := n.producer(done, stream)

	// fan-out
	var wg sync.WaitGroup
	wg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go n.consumer(tweets, &wg)
	}

	wg.Wait()
	fmt.Printf("Concurrent process took %s\n", time.Since(start))
}

func (n *ConcurrentProcessor) producer(done <-chan struct{}, stream Stream) <-chan *Tweet {
	tweets := make(chan *Tweet)
	go func() {
		defer close(tweets)
		for {
			select {
			case <-done:
			default:
				tweet, err := stream.Next()
				if err == ErrEOF {
					return
				}
				tweets <- tweet
			}
		}
	}()
	return tweets
}

func (n *ConcurrentProcessor) consumer(tweets <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}
