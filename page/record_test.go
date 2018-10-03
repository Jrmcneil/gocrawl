package page

import (
	"gocrawl/job"
	"testing"
)

func TestRecord(t *testing.T) {
	t.Run("Passes a new page to the output channel", func(t *testing.T) {
		in := make(chan job.Job)
		out := make(chan job.Job)
		quit := make(chan bool, 1)
		record := NewRecord(in, out, quit)
		page := NewPage("www.monzo.com")

		record.Start()
		in <- page
		savedPage := <-out

		if savedPage.Address() != page.address {
			t.Errorf("got: %s, want: %s", savedPage.Address(), page.address)
		}
	})

	t.Run("Does not pass a visited page to the output channel", func(t *testing.T) {
		in := make(chan job.Job)
		out := make(chan job.Job)
		quit := make(chan bool, 1)
		record := NewRecord(in, out, quit)
		page := NewPage("www.monzo.com")

		record.Start()
		go func() {
			in <- page
			in <- page
		}()

		<-out

		if len(out) != 0 {
			t.Errorf("got: %d, want: %d", len(out), 0)
		}
	})

	t.Run("Increases the record with each new page", func(t *testing.T) {
		in := make(chan job.Job)
		out := make(chan job.Job)
		quit := make(chan bool, 1)
		record := NewRecord(in, out, quit)
		page1 := NewPage("www.monzo.com")
		page2 := NewPage("www.monzo.com/contact")

		record.Start()
		go func() {
			in <- page1
			in <- page1
			in <- page2
		}()

		<-out
		<-out

		if len(record.record) != 2 {
			t.Errorf("got: %d, want: %d", len(record.record), 2)
		}
	})
}
