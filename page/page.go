package page

import (
	"regexp"
	"sync"
)

type Page struct {
	address string
	links   []*Page
	Done    chan bool
}

func (page *Page) Address() string {
	return page.address
}

func (page *Page) Links() []*Page {
	return page.links
}

func (page *Page) ParseLinks(html string) {
	matches := findMatches(page.address, html)
	page.links = setLinks(stripURL(page.address), matches)
	page.watchLinks()
}

func (page *Page) watchLinks() {
	go func(page *Page) {
		var wg sync.WaitGroup
		wg.Add(len(page.links))
		for _, link := range page.links {
			<-link.Done
			wg.Done()
		}
		wg.Wait()
		page.Done <- true
	}(page)
}

func setLinks(domain string, matches [][]string) []*Page {
	links := make([]*Page, len(matches))
	for index, match := range matches {
		links[index] = NewPage(domain + match[1])
	}
	return links
}

func findMatches(address string, html string) [][]string {
	regex := regexp.MustCompile(`<a[ ]+href=\".*` + escapedDomain(address) + `(/.*)\".*</a>`)
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
	page.Done = make(chan bool, 1)
	return page
}
