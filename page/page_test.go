package page

import (
    "gocrawl/job"
    "testing"
)

func TestPage(t *testing.T) {
	const htmlString = `
        <a href="https://www.monzo.com/accounts"></a>
		<a target="_self" href="http://monzo.co.uk/advice"></a>
		<a href="http://www.monzo.com/banking/loans.php"></a>
		<a href="http://www.monzo.com/"> Some words</a>
        <a href="www.monzo.co.uk/"></a>
        <a href="www.google.com/something/"></a>
		<a href="https://www.monzo.com/contact/london"> </a>
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

	t.Run("Build extracts links from html string using domain as address", func(t *testing.T) {

		page := NewPage("monzo.co.uk")
		page.Build(htmlString)

		expected := map[string]bool{
			"monzo.co.uk/advice": true,
			"monzo.co.uk/":       true}

		assertCorrectLinks(t, page.links, expected)
	})

	t.Run("Build extracts links from html string using URL as address", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)

		expected := map[string]bool{
			"www.monzo.com/accounts":          true,
			"www.monzo.com/banking/loans.php": true,
			"www.monzo.com/contact/london":    true,
			"www.monzo.com/":                  true}

		assertCorrectLinks(t, page.links, expected)
	})

	t.Run("ResetLinks removes references to linked Jobs", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)
		page.ResetLinks()

		if len(page.links) != 0 {
			t.Errorf("got: %v, want: %v", len(page.links), 0)
		}
	})

	t.Run("Pages are not ready when created", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")

		if len(page.Ready()) != 0 {
			t.Errorf("got: %d, want: %d", len(page.Ready()), 0)
		}
	})

	t.Run("Pages are ready when built", func(t *testing.T) {
		page := NewPage("https://www.monzo.com/contact/london/")
		page.Build(htmlString)

		if len(page.Ready()) != 1 {
			t.Errorf("got: %d, want: %d", len(page.Ready()), 1)
		}
	})
}
