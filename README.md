GoCrawl
=====


Description
----

A simple website crawler. Pass it a url and it will find and print a sitemap of all the links it can reach by crawling from there.

How to run
----

Install the dependencies with ```go get ./...```

Install the crawler globally with ```go install``` or locally with ```go build```. 

To view a sitemap, run the executable with the root url of the site you want to crawl. You can also configure the number of workers, queue sizes and HTTP request rate limit instead of using the defaults.

```bash
gocrawl -w 15 -q 30 -l 300 http://www.example.com
```

for help check ```gocrawl -h```

Tuning
----

The crawler uses a [worker pool](https://gobyexample.com/worker-pools) to process pages they receive on a queue from a central record. Workers fetch links on the pages and send those as pages on another queue back to the record for checking. 

Increasing the number of workers increases the number of pages that can be processed concurrently but can require a larger queue depending on the structure of the site.

The program will exit with a warning if workers become overloaded and cannot send pages to the record.


To Do
---

- Expand link parsing to cover relative links. Currently only absolute links will be processed as new pages.
- Under the right circumstances the record could increase until it overflows the heap - i.e. a huge website with a manageable number of new links per page. Consider setting an upper bound to the record size
- Refactor tests