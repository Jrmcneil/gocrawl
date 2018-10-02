package page

import "testing"

func TestSync(t *testing.T) {
    t.Run("Passes a new page to the output channel", func(t *testing.T) {
       in := make(chan *Page)
       out := make(chan *Page)
       quit := make(chan bool, 1)
       sync := NewSync(in, out, quit)
       page := NewPage("www.monzo.com")

       sync.Start()
       in <- page
       savedPage := <- out

       if savedPage.address != page.address {
           t.Errorf("got: %s, want: %s", savedPage.address, page.address)
       }
    })

    t.Run("Does not pass a visited page to the output channel", func(t *testing.T) {
        in := make(chan *Page)
        out := make(chan *Page)
        quit := make(chan bool, 1)
        sync := NewSync(in, out, quit)
        page := NewPage("www.monzo.com")

        sync.Start()
        go func() {
            in <- page
            in <- page
        }()

        <- out


        if len(out) != 0 {
            t.Errorf("got: %d, want: %d", len(out), 0)
        }
    })

    t.Run("Increases the record with each new page", func(t *testing.T) {
        in := make(chan *Page)
        out := make(chan *Page)
        quit := make(chan bool, 1)
        sync := NewSync(in, out, quit)
        page1 := NewPage("www.monzo.com")
        page2 := NewPage("www.monzo.com/contact")

        sync.Start()
        go func() {
            in <- page1
            in <- page1
            in <- page2
        }()

        <- out
        <- out

        if len(sync.record) != 2 {
            t.Errorf("got: %d, want: %d", len(sync.record), 2)
        }
    })
}
