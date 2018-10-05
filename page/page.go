package page

import (
	"gocrawl/job"
	"regexp"
)

type Page struct {
	address string
	links   []job.Job
	ready   chan bool
}

func (page *Page) Ready() chan bool {
	return page.ready
}

func (page *Page) Address() string {
	return page.address
}

func (page *Page) Links() []job.Job {
	return page.links
}

func (page *Page) ResetLinks() {
	page.links = make([]job.Job, 0)
}

func (page *Page) Build(html string) {
	matches := findMatches(page.address, html)
	page.links = setLinks(stripURL(page.address), matches)
	page.ready <- true
}

func setLinks(domain string, matches [][]string) []job.Job {
	links := make([]job.Job, len(matches))
	for index, match := range matches {
		links[index] = NewPage(domain + match[1])
	}
	return links
}

func findMatches(address string, html string) [][]string {
	regex := regexp.MustCompile(`<a[ ]+.*href=\".*` + escapedDomain(address) + `(\/.*)\".*>.*<\/a>`)
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
	return page
}
