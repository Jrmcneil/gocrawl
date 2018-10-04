package gocrawl

import (
	"flag"
    "fmt"
    "gocrawl/client"
	"gocrawl/crawler"
)

var (
	workers = flag.Int("w", 10, "Number of job workers to initialise")
)

func main() {
	flag.Parse()
	url := flag.Arg(0)
	c := crawler.NewCrawler(*workers, client.NewClient)
	result := c.Crawl(url)
	fmt.Println(result)
}
