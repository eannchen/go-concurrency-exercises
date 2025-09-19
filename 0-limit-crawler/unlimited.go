package main

import (
	"fmt"
	"sync"
)

// UnlimitedCrawler is a Crawler that fetches pages as fast as possible.
type UnlimitedCrawler struct{}

// Crawl starts the crawling process and waits for it to complete.
func (c *UnlimitedCrawler) Crawl(url string, depth int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go c.crawlRecursive(url, depth, &wg)
	wg.Wait()
}

// crawlRecursive is the internal, recursive crawling logic.
func (c *UnlimitedCrawler) crawlRecursive(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

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
