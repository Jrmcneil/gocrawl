package page

import "regexp"

type Page struct {
	address string
	links   []string
}

func (page *Page) ParseLinks(html string) {
    matches := findMatches(page.address, html)
    page.links = setLinks(matches)
}

func setLinks(matches [][]string) []string {
    links := make([]string, len(matches))
    for index, match := range matches {
        links[index] = match[1]
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
	regex := regexp.MustCompile(`(?:http://|https://|)(?:.+\.|)(.+\.[^/]+)\/?.*`)
	match := regex.FindStringSubmatch(address)
	return match[1]
}

func NewPage(address string) *Page {
	return &Page{address: address}
}
