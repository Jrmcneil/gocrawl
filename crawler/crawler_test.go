package crawler

import (
	"fmt"
	"gocrawl/client"
	"testing"
	"time"
)

func TestCrawler(t *testing.T) {

	htmlMap := make(map[string]string, 3)

	root := `
<html>
    <head>
        <title>Monzo.com</title>
    </head>
    <body>
        <div>
 <a href="https://www.monzo.com/about">About Us</a>
           
        </div>

       <a target="_self" href="https://www.monzo.com/contact">Contact Us</a> 

        <a href="https://www.monzo.com/reviews">Reviews</a>
    </body>
</html>
`

	about := `
<html>
    <head>
        <title>About</title>
    </head>
    <body>
        <div>
            <a target="_self" href="https://www.monzo.com/">Home</a>
        </div>

        <div>
            <a target="_self" href="https://www.monzo.com/contact">Contact Us</a>
        </div>
    </body>
</html>
`
	contact := `
<html>
    <head>
        <title>Contacts</title>
    </head>
    <body>
        Find us at info@monzo.com
    </body>
</html>
`

	reviews := `
<html>
    <head>
        <title>Reviews</title>
    </head>
    <body>
        Check us out at <a href="https://www.trustpilot.com/monzo">Trustpilot</a>
    </body>
</html>
`

	htmlMap["https://www.monzo.com/"] = root
	htmlMap["www.monzo.com/about"] = about
	htmlMap["www.monzo.com/contact"] = contact
	htmlMap["www.monzo.com/reviews"] = reviews

	t.Run("Crawl returns a sitemap", func(t *testing.T) {

		c := NewCrawler(1, 2, testClientBuilder(htmlMap), 1)
		result := c.Crawl("https://www.monzo.com/")

		fmt.Println(result)
	})
}

type TestClient struct {
	htmlMap map[string]string
}

func testClientBuilder(htmlMap map[string]string) func(limiter <-chan time.Time) client.HttpClient {
	return func(limiter <-chan time.Time) client.HttpClient {
		return &TestClient{htmlMap: htmlMap}
	}
}

func (client *TestClient) Get(address string) (string, error) {
	return client.htmlMap[address], nil
}
