package page

import (
	"testing"
)

func TestPage(t *testing.T) {
	const htmlString = `
        <a href="https://www.monzo.com/accounts"</a>
        <a href="monzo.com/support"</a>
		<a href="http://monzo.com/advice"</a>
		<a href="http://www.monzo.com/banking/loans"</a>
	`

	t.Run("New Page extracts links from html string", func(t *testing.T) {

		page := NewPage("monzo.com", htmlString)

		expected := map[string]bool{
			"/accounts":      true,
			"/support":       true,
			"/advice":        true,
			"/banking/loans": true}

		for _, link := range page.links {
			if expected[link] != true {
				t.Errorf("'%s' was not present in links", link)
			}
		}

		if len(page.links) != len(expected) {
			t.Errorf("expected %d links got %d ", 2, len(page.links))
		}
	})

	t.Run("Passing a URL matches links where only domain is specified", func(t *testing.T) {
		page := NewPage("https://www.monzo.com", htmlString)

		expected := map[string]bool{
			"/accounts":      true,
			"/support":       true,
			"/advice":        true,
			"/banking/loans": true}

		for _, link := range page.links {
			if expected[link] != true {
				t.Errorf("'%s' was not present in links", link)
			}
		}

		if len(page.links) != len(expected) {
			t.Errorf("expected %d links got %d ", 2, len(page.links))
		}
	})

}
