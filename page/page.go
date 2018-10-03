package page

import (
    "gocrawl/job"
    "regexp"
    "sync"
)

type Page struct {
	address string
	links   []job.Job
	done    chan bool
	ready   chan bool
}

func (page *Page) Ready() chan bool {
	return page.ready
}

func (page *Page) Done() chan bool {
	return page.done
}

func (page *Page) Address() string {
	return page.address
}

func (page *Page) Links() []job.Job {
	return page.links
}

func (page *Page) Build(html string) {
	matches := findMatches(page.address, html)
	page.links = setLinks(stripURL(page.address), matches)
	page.ready <- true
	page.watchLinks()
}

func (page *Page) Close() {
	page.done <- true
}

func (page *Page) watchLinks() {
	go func(page *Page) {
		var wg sync.WaitGroup
		wg.Add(len(page.links))
		for _, link := range page.links {
			<-link.Done()
			wg.Done()
		}
		wg.Wait()
		page.Close()
	}(page)
}

func setLinks(domain string, matches [][]string) []job.Job {
	links := make([]job.Job, len(matches))
	for index, match := range matches {
		links[index] = NewPage(domain + match[1])
	}
	return links
}

func findMatches(address string, html string) [][]string {
	regex := regexp.MustCompile(`<a[ ]+.*href=\".*` + escapedDomain(address) + `(\/.*)\".*<\/a>`)
	return regex.FindAllStringSubmatch(html, -1)
}

func escapedDomain(address string) string {
	domain := stripURL(address)
	return regexp.QuoteMeta(domain)
}

func stripURL(address string) string {
	regex := regexp.MustCompile(`(?:^https:\/\/|^http:\/\/|^)([^\/]+)\/?.*$`)
	match := regex.FindStringSubmatch(address)
	return match[1]
}

func NewPage(address string) *Page {
	page := new(Page)
	page.address = address
	page.ready = make(chan bool, 1)
	page.done = make(chan bool, 1)
    return page
}
