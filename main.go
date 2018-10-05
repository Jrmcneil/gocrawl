package gocrawl

import (
	"flag"
	"fmt"
	"gocrawl/client"
	"gocrawl/crawler"
)

var (
	workers   = flag.Int("w", 10, "Number of job workers to initialise")
	queueSize = flag.Int("q", 10, "Size of input and output channels to the page record")
	rateLimit = flag.Int("l", 50, "Time (in milliseconds) between HTTP requests")
)

func main() {
	flag.Parse()
	url := flag.Arg(0)
	c := crawler.NewCrawler(*workers, *queueSize, client.NewClientBuilder(), *rateLimit)
	result := c.Crawl(url)
	fmt.Println(result)
}
