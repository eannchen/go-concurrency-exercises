package main

import (
	"testing"
	"time"
)

// TestLimitedCrawler verifies that the rate-limited crawler does not fetch
// pages faster than one per second.
func TestLimitedCrawler(t *testing.T) {
	fetchSig := fetchSignalInstance()

	// This goroutine acts as a high-resolution stopwatch. It checks the
	// time between consecutive fetch signals.
	go func() {
		lastFetchTime := time.Time{} // Initialize to zero time
		for range fetchSig {
			if !lastFetchTime.IsZero() {
				// Check if the duration is less than 1 second (with a small tolerance).
				if time.Since(lastFetchTime) < 950*time.Millisecond {
					t.Errorf("Two crawls were executed less than 1 second apart.")
				}
			}
			lastFetchTime = time.Now()
		}
	}()

	// Instantiate and run the specific crawler we want to test.
	crawler := &LimitedCrawler{}
	crawler.Crawl("http://golang.org/", 4)

	// Close the signal channel to allow the test to finish cleanly.
	close(fetchSignal)
}
