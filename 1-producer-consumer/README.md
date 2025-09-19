# Producer-Consumer Scenario

The producer reads in tweets from a mockstream and a consumer is processing the data to find out whether someone has tweeted about golang or not.

## Expected result
```
--- Running Sequential Processor ---
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
vampirewalk666  tweets about golang
Sequential process took 3.58755825s

--- Running Concurrent Processor ---
davecheney      tweets about golang
beertocode      does not tweet about golang
ironzeb         tweets about golang
beertocode      tweets about golang
vampirewalk666  tweets about golang
Concurrent process took 1.936573625s
```