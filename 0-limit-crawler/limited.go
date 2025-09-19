package main

import (
	"fmt"
	"sync"
	"time"
)

// rateLimiter is the shared ticker that enforces the 1 fetch/sec rule.
var rateLimiter = time.Tick(time.Second)

// LimitedCrawler is a Crawler that fetches at most one page per second.
type LimitedCrawler struct{}

// Crawl starts the rate-limited crawling process and waits for it to complete.
func (c *LimitedCrawler) Crawl(url string, depth int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go c.crawlRecursive(url, depth, &wg)
	wg.Wait()
}

// crawlRecursive is the internal logic that includes the rate-limiting step.
func (c *LimitedCrawler) crawlRecursive(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	<-rateLimiter // This line blocks until the next tick from the global ticker.

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		go c.crawlRecursive(u, depth-1, wg)
	}
}
