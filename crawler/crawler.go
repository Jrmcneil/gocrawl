package crawler

import (
	"gocrawl/client"
	"gocrawl/job"
	"gocrawl/page"
	"gocrawl/record"
	"gocrawl/sitemap"
	"gocrawl/worker"
	"log"
	"os"
	"time"
)

type Crawler struct {
	workers  []*worker.Worker
	record   *record.Record
	pipeline chan job.Job
	queue    chan job.Job
	tree     *sitemap.Sitemap
}

func (c *Crawler) Crawl(url string) string {
	defer func() { c.stop() }()
	root := page.NewPage(url)
	c.start(root, url)
	output := <-c.tree.Result
	return output
}

func (c *Crawler) stop() {
	c.record.Stop()
	c.stopWorkers()
}

func (c *Crawler) start(root job.Job, url string) {
	c.startWorkers(url)
	go c.record.Start()
	go c.tree.Build(root)
	c.pipeline <- root
}

func (c *Crawler) stopWorkers() {
	workers := len(c.workers)
	for i := 0; i < workers; i++ {
		c.workers[i].Stop()
	}
}

func (c *Crawler) startWorkers(url string) {
	workers := len(c.workers)
	c.watchWorkers(url)
	for i := 0; i < workers; i++ {
		go c.workers[i].Start()
	}
}

func (c *Crawler) watchWorkers(url string) {
	for i := range c.workers {
		go func(worker *worker.Worker) {
			<-worker.Overload()
			log.Printf("WARNING: Workers overloaded. Could not complete sitemap. Increase queue size to fully process %s", url)
			c.stop()
			os.Exit(1)
		}(c.workers[i])
	}
}

func NewCrawler(workers int, queueSize int, clientBuilder client.HttpClientBuilder, rate int) *Crawler {
	crawler := new(Crawler)
	crawler.queue = make(chan job.Job, queueSize)
	crawler.pipeline = make(chan job.Job, queueSize)
	crawler.record = record.NewRecord(crawler.pipeline, crawler.queue, make(chan bool, 1))
	crawler.workers = make([]*worker.Worker, workers)
	crawler.tree = &sitemap.Sitemap{Result: make(chan string, 1)}
	limiter := time.Tick(time.Duration(rate) * time.Millisecond)
	for i := 0; i < workers; i++ {
		crawler.workers[i] = worker.NewWorker(crawler.queue, crawler.pipeline, clientBuilder(limiter))
	}
	return crawler
}
