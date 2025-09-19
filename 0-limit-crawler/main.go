package main

import (
	"fmt"
	"time"
)

// Crawler defines the interface for crawling a website.
type Crawler interface {
	// Crawl starts crawling from a given URL up to a certain depth.
	Crawl(url string, depth int)
}

func main() {
	crawlers := map[string]Crawler{
		"Unlimited": &UnlimitedCrawler{},
		"Limited":   &LimitedCrawler{},
	}

	for name, crawler := range crawlers {
		fmt.Printf("--- Running %s Crawler ---\n", name)
		start := time.Now()
		crawler.Crawl("http://golang.org/", 4)
		fmt.Printf("--- %s Crawler took %s ---\n\n", name, time.Since(start))
	}
}
