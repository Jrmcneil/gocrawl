package page

import (
	"gocrawl/job"
	"testing"
)

func TestPage(t *testing.T) {
	const htmlString = `
        <a href="https://www.monzo.com/accounts"</a>
		<a href="http://monzo.co.uk/advice"</a>
		<a href="http://www.monzo.com/banking/loans.php"</a>
		<a href="http://www.monzo.com/"</a>
        <a href="www.monzo.co.uk/"</a>
		<a href="https://www.monzo.com/contact/london" </a>
	`

	assertCorrectLinks := func(t *testing.T, links []job.Job, expected map[string]bool) {
		t.Helper()
		for _, link := range links {
			if expected[link.Address()] != true {
				t.Errorf("'%s' was not expected in links", link.Address())
			}
		}

		if len(links) != len(expected) {
			t.Errorf("expected %d links got %d ", len(expected), len(links))
		}
	}

	t.Run("New Page extracts links from html string using domain as address", func(t *testing.T) {

		page := NewPage("monzo.co.uk")
		page.Build(htmlString)

		expected := map[string]bool{
			"monzo.co.uk/advice": true,
			"monzo.co.uk/":       true}

		assertCorrectLinks(t, page.links, expected)
	})

	t.Run("New Page extracts links from html string using URL as address", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)

		expected := map[string]bool{
			"www.monzo.com/accounts":          true,
			"www.monzo.com/banking/loans.php": true,
			"www.monzo.com/contact/london":    true,
			"www.monzo.com/":                  true}

		assertCorrectLinks(t, page.links, expected)
	})

	t.Run("Pages are not closed by default", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)

		done := len(page.Done()) != 0

		if done != false {
			t.Errorf("got: %t, want: %t", done, false)
		}
	})

	t.Run("Page closes once its links are closed", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)

		for _, link := range page.links {
			link.Close()
		}

		done := <-page.Done()
		if done != true {
			t.Errorf("got: %t, want: %t", done, true)
		}
	})

}
