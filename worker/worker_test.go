package worker

import (
    "testing"
)

func TestWorker(t *testing.T) {

	t.Run("Start gets the html for the given address", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan Job)
		record := make(chan Job)

		worker := NewWorker(queue, record, client)

		address := "www.monzo.com"
        p := &TestJob{calls: make(map[string]int, 4)}
		p.address = address

		worker.Start()
		worker.queue <- p

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
		queue := make(chan Job)
		record := make(chan Job)

		worker := NewWorker(queue, record, client)

        address := "www.monzo.com"
        p := &TestJob{calls: make(map[string]int, 4)}
        p.address = address
        link := &TestJob{calls: make(map[string]int, 4)}
        p.links = []Job{link}

		worker.Start()

		queue <- p
		sent := <-record

		if sent != link {
			t.Errorf("got: %v, want: %v", sent, p)
		}
	})

	t.Run("Start marks the page as done if there is an error and doesn't add to record", func(t *testing.T) {
		client := new(TestClient)
		client.args = make([]string, 1)
		queue := make(chan Job)
		record := make(chan Job)

		worker := NewWorker(queue, record, client)

		address := "www.monzo.com"
        p := &TestJob{calls: make(map[string]int, 4)}
        p.address = address

        client.err = &TestError{}

        worker.Start()
		worker.queue <- p

		callsToClose := p.calls["Close"]

		if callsToClose != 1 {
			t.Errorf("got: %d, want: %d", callsToClose, 1)
		}

		if len(record) != 0 {
			t.Errorf("got: %d, want: %d", len(record), 0)
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
    links []Job
    calls map[string]int
    args  []string
    address string
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

func (job *TestJob) Links() []Job {
    job.calls["Links"] = job.calls["Links"] + 1
    return job.links
}

type TestError struct{}

func (err *TestError) Error() string {
    return "Test error"
}
