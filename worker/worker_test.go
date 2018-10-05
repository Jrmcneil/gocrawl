package worker

import (
	"gocrawl/job"
	"testing"
)

func TestWorker(t *testing.T) {

	t.Run("Start gets the html for the given address", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan job.Job)
		record := make(chan job.Job)

		worker := NewWorker(queue, record, client)

		address := "www.monzo.com"
		p := &TestJob{calls: make(map[string]int, 4)}
		p.address = address

		go worker.Start()
		queue <- p

		if client.calls != 1 {
			t.Errorf("got: %d, want: %d", client.calls, 1)
		}

		if client.args[1] != address {
			t.Errorf("got: %s, want: %s", client.args[1], address)
		}
	})

	t.Run("Start adds the page's links to the record", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan job.Job, 2)
		record := make(chan job.Job, 1)

		worker := NewWorker(queue, record, client)

		address := "www.monzo.com"
		p := &TestJob{calls: make(map[string]int, 4)}
		p.address = address
		link := &TestJob{calls: make(map[string]int, 4)}
		p.links = []job.Job{link}

		go worker.Start()

		queue <- p
		sent := <-record

		if sent != link {
			t.Errorf("got: %v, want: %v", sent, p)
		}
	})

	t.Run("Start marks the page as ready if there is an error and doesn't add to record", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan job.Job, 1)
		record := make(chan job.Job, 1)
		ready := make(chan bool, 1)

		worker := NewWorker(queue, record, client)

		p := &TestJob{calls: make(map[string]int, 4), ready: ready}
		p.address = "www.monzo.com"

		client.err = &TestError{}

		go worker.Start()
		queue <- p

		r := <-ready

		if r != true {
			t.Errorf("got: %d, want: %d", r, true)
		}

		if len(record) != 0 {
			t.Errorf("got: %d, want: %d", len(record), 0)
		}
	})

	t.Run("Overloads if the record is full", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan job.Job, 1)
		record := make(chan job.Job, 1)

		worker := NewWorker(queue, record, client)

		p1 := &TestJob{calls: make(map[string]int, 4)}
		p1.address = "www.monzo.com"
		link1 := &TestJob{calls: make(map[string]int, 4)}
		p1.links = []job.Job{link1}

		p2 := &TestJob{calls: make(map[string]int, 4)}
		p2.address = "www.monzo.com/about"
		link2 := &TestJob{calls: make(map[string]int, 4)}
		p2.links = []job.Job{link2}

		go worker.Start()
		queue <- p1
		queue <- p2

		overload := <-worker.overload

		if overload != true {
			t.Errorf("got: %d, want: %d", overload, true)
		}
	})
}

type TestClient struct {
	calls int
	args  []string
	err   error
}

func (client *TestClient) Get(address string) (string, error) {
	client.calls++
	client.args = append(client.args, address)
	return address, client.err
}

type TestJob struct {
	links   []job.Job
	calls   map[string]int
	args    []string
	address string
	ready   chan bool
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

func (job *TestJob) Ready() chan bool {
	job.calls["Ready"] = job.calls["Ready"] + 1
	return job.ready
}
func (job *TestJob) ResetLinks() {}

type TestError struct{}

func (err *TestError) Error() string {
	return "Test error"
}
