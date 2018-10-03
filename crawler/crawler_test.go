package crawler

import (
    "fmt"
    "gocrawl/client"
    "testing"
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
            <a target="_self" href="https://www.monzo.com/contact">Contact Us</a>
        </div>

        <a href="https://www.monzo.com/about">About Us</a>
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


htmlMap["https://www.monzo.com/"] = root
htmlMap["www.monzo.com/about"] = about
htmlMap["www.monzo.com/contact"] = contact

    t.Run("Crawl returns a sitemap", func(t *testing.T) {

        c := NewCrawler(3, testClientBuilder(htmlMap))
        result := c.Crawl("https://www.monzo.com/")

        fmt.Println(result)
    })
}

type TestClient struct {
    htmlMap map[string]string
}

func testClientBuilder(htmlMap map[string]string) func() client.HttpClient {
    return func() client.HttpClient {
        return &TestClient{htmlMap: htmlMap}
    }
}


func (client *TestClient) Get(address string) (string, error) {
    return client.htmlMap[address], nil
}

