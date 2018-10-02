package page

import (
	"testing"
)

func TestPage(t *testing.T) {
	const htmlString = `
        <a href="https://www.monzo.com/accounts"</a>
		<a href="http://monzo.com/advice"</a>
		<a href="http://www.monzo.com/banking/loans.php"</a>
		<a href="http://www.monzo.com/"</a>
		<a href="https://www.monzo.com/contact/london" </a>
	`

	assertCorrectLinks := func(t *testing.T, links []string, expected map[string]bool) {
		t.Helper()
		for _, link := range links {
			if expected[link] != true {
				t.Errorf("'%s' was not expected in links", link)
			}
		}

		if len(links) != len(expected) {
			t.Errorf("expected %d links got %d ", len(expected), len(links))
		}
	}

	t.Run("New Page extracts links from html string using domain as address", func(t *testing.T) {

		page := NewPage("monzo.com")
		page.ParseLinks(htmlString)

		expected := map[string]bool{
			"/accounts":          true,
			"/advice":            true,
			"/banking/loans.php": true,
			"/contact/london":    true,
			"/":                  true}

		assertCorrectLinks(t, page.links, expected)
	})

	t.Run("New Page extracts links from html string using URL as address", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.ParseLinks(htmlString)

		expected := map[string]bool{
			"/accounts":          true,
			"/advice":            true,
			"/banking/loans.php": true,
			"/contact/london":    true,
			"/":                  true}

		assertCorrectLinks(t, page.links, expected)
	})

}
