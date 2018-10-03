package record

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
		page := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com"}

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
		page := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com"}

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

	t.Run("Closes a visited page", func(t *testing.T) {
		in := make(chan job.Job)
		out := make(chan job.Job)
		quit := make(chan bool, 1)
		record := NewRecord(in, out, quit)
		page1 := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com"}
		page2 := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com"}

		record.Start()
		go func() {
			in <- page1
			in <- page2
		}()

		<-out

		if len(page2.Done()) != 0 {
			t.Errorf("got: %d, want: %d", len(page2.Done()), 0)
		}
	})

	t.Run("Increases the record with each new page", func(t *testing.T) {
		in := make(chan job.Job)
		out := make(chan job.Job)
		quit := make(chan bool, 1)
		record := NewRecord(in, out, quit)
		page1 := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com"}
		page2 := &TestJob{calls: make(map[string]int, 4), address: "www.monzo.com/contact"}

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

type TestJob struct {
	links   []job.Job
	calls   map[string]int
	args    []string
	address string
	done    chan bool
}

func (job *TestJob) Address() string {
	job.calls["Address"] = job.calls["Address"] + 1
	return job.address
}

func (job *TestJob) Close() {
	job.calls["Close"] = job.calls["Close"] + 1
}

func (job *TestJob) Build(str string) {
	job.args = append(job.args, str)
	job.calls["Build"] = job.calls["Build"] + 1
}

func (job *TestJob) Links() []job.Job {
	job.calls["Links"] = job.calls["Links"] + 1
	return job.links
}

func (job *TestJob) Done() chan bool {
	job.calls["Done"] = job.calls["Done"] + 1
	return job.done
}

func (job *TestJob) Ready() chan bool {
	job.calls["Ready"] = job.calls["Ready"] + 1
	return make(chan bool, 1)
}
