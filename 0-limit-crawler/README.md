# Limit your crawler

Given is a crawler (modified from the Go tour) that requests pages
excessively. However, we don't want to burden the webserver too
much. The task is to change the code to limit the crawler to at most
one page per second, while maintaining concurrency (in other words,
Crawl() must be called concurrently)

## Hint

This exercise can be solved in 3 lines only. If you can't do
it, have a look at this:
https://go.dev/wiki/RateLimiting

## Expected result

```
--- Running Unlimited Crawler ---
found: http://golang.org/ "The Go Programming Language"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/pkg/os/ "Package os"
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/ "The Go Programming Language"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/ "Packages"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/fmt/ "Package fmt"
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/ "The Go Programming Language"
found: http://golang.org/ "The Go Programming Language"
--- Unlimited Crawler took 262.458µs ---

--- Running Limited Crawler ---
found: http://golang.org/ "The Go Programming Language"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/pkg/os/ "Package os"
found: http://golang.org/ "The Go Programming Language"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/fmt/ "Package fmt"
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/ "The Go Programming Language"
not found: http://golang.org/cmd/
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/pkg/ "Packages"
found: http://golang.org/ "The Go Programming Language"
--- Limited Crawler took 13.000772458s ---
```
