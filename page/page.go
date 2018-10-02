package page

import "regexp"

type Page struct {
	domain string
	html   string
	links  []string
}

func parseLinks(domain string, html string) []string {
	escapedDomain := regexp.QuoteMeta(domain)
	regex := regexp.MustCompile(`<a[ ]+href=\".*` + escapedDomain + `(/.*)\".*</a>`)
	matches := regex.FindAllStringSubmatch(html, -1)
	links := make([]string, len(matches))
	for index, match := range matches {
		links[index] = match[1]
	}

	return links
}

func stripURL(address string) string {
	regex := regexp.MustCompile(`(?:http://|https://|)(?:.+\.|)(.+\.[^/]+)\/?.*`)
	match := regex.FindStringSubmatch(address)
	return match[1]
}

func NewPage(address string, html string) *Page {
	page := new(Page)
	page.html = html
	page.domain = stripURL(address)
	page.links = parseLinks(page.domain, html)
	return page
}
