package page

import (
    "testing"
)

func TestPageHasLinks(t *testing.T) {
    htmlString := `<html>
    <body>
    <a href="https://www.monzo.com/accounts"</a>
    <a href="monzo.com/support"</a>
    </body>
    </html>`
    page := NewPage("monzo.com", htmlString)

    expected := map[string]bool{"/accounts": true, "/support": true}

    for _, link := range page.links {
        if expected[link] != true {
            t.Errorf("'%s' was not present in links", link)
        }
    }

    if len(page.links) != 2 {
        t.Errorf("expected %d links got %d ", 2, len(page.links))
    }
}