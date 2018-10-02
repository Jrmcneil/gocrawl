package page

import "regexp"

type Page struct {
	address string
	links  []string
}

func (page *Page) ParseLinks(html string) {
    domain := stripURL(page.address)
	escapedDomain := regexp.QuoteMeta(domain)
	regex := regexp.MustCompile(`<a[ ]+href=\".*` + escapedDomain + `(/.*)\".*</a>`)
	matches := regex.FindAllStringSubmatch(html, -1)
	links := make([]string, len(matches))
	for index, match := range matches {
		links[index] = match[1]
	}

	page.links = links
}

func stripURL(address string) string {
	regex := regexp.MustCompile(`(?:http://|https://|)(?:.+\.|)(.+\.[^/]+)\/?.*`)
	match := regex.FindStringSubmatch(address)
	return match[1]
}

func NewPage(address string) *Page {
	return &Page{address: address}
}
