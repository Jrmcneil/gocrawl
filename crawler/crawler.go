package crawler

import (
    "gocrawl/client"
    "gocrawl/job"
    "gocrawl/page"
    "gocrawl/record"
    "gocrawl/sitemap"
    "gocrawl/worker"
)

type Crawler struct {
    workers []*worker.Worker
    record *record.Record
    pipeline chan job.Job
    queue chan job.Job
    tree *sitemap.Sitemap
}

func (c *Crawler) Crawl(url string) string {
    defer c.stop()
    root := page.NewPage(url)
    c.start(root)
    return <- c.tree.Result
}

func (c *Crawler) stop() {
    c.record.Stop()
    c.stopWorkers()
}

func (c *Crawler) start(root job.Job) {
    c.startWorkers()
    go c.record.Start()
    go c.tree.Build(root)
    c.pipeline <- root
}

func (c* Crawler) stopWorkers() {
    workers := len(c.workers)
    for i := 0; i < workers; i++ {
        c.workers[i].Stop()
    }
}

func (c *Crawler) startWorkers() {
    workers := len(c.workers)
    for i := 0; i < workers; i++ {
        go c.workers[i].Start()
    }
}

func NewCrawler(workers int, queueSize int, clientBuilder client.HttpClientBuilder) *Crawler {
   crawler := new(Crawler)
   crawler.queue = make(chan job.Job, queueSize)
   crawler.pipeline = make(chan job.Job, queueSize)
   crawler.record = record.NewRecord(crawler.pipeline, crawler.queue, make(chan bool, 1))
   crawler.workers = make([]*worker.Worker, workers)
   crawler.tree = &sitemap.Sitemap{Result: make(chan string, 1)}
   for i := 0; i < workers; i++ {
       crawler.workers[i] = worker.NewWorker(crawler.queue, crawler.pipeline, clientBuilder())
   }
   return crawler
}
