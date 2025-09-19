package main

import (
	"fmt"
	"time"
)

// SequentialProcessor handles processing tweets one by one.
type SequentialProcessor struct{}

// Process runs the entire sequential pipeline: producer -> consumer.
func (n *SequentialProcessor) Process() {
	start := time.Now()
	stream := GetMockStream()

	tweets := n.producer(stream)
	n.consumer(tweets)

	fmt.Printf("Sequential process took %s\n", time.Since(start))
}

func (n *SequentialProcessor) producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}
		tweets = append(tweets, tweet)
	}
}

func (n *SequentialProcessor) consumer(tweets []*Tweet) {
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}
