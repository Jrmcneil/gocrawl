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
    root := page.NewPage(url)
    c.startWorkers()
    c.record.Start()
    c.pipeline <- root
    c.tree.Build(root)
    return <- c.tree.Result
}

func (c *Crawler) startWorkers() {
    workers := len(c.workers)
    for i := 0; i < workers; i++ {
        c.workers[i].Start()
    }
}

func NewCrawler(workers int, clientBuilder client.HttpClientBuilder) *Crawler {
   crawler := new(Crawler)
   crawler.queue = make(chan job.Job, 10)
   crawler.pipeline = make(chan job.Job, 10)
   crawler.record = record.NewRecord(crawler.pipeline, crawler.queue, make(chan bool, 1))
   crawler.workers = make([]*worker.Worker, workers)
   crawler.tree = &sitemap.Sitemap{Result: make(chan string, 1)}
   for i := 0; i < workers; i++ {
       crawler.workers[i] = worker.NewWorker(crawler.queue, crawler.pipeline, clientBuilder())
   }
   return crawler
}
