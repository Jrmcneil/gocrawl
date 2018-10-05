package worker

import (
	"gocrawl/client"
	"gocrawl/job"
)

type Worker struct {
	client   client.HttpClient
	queue    <-chan job.Job
	record   chan<- job.Job
	quit     chan bool
	overload chan bool
}

func (worker *Worker) Start() {
	for {
		select {
		case p := <-worker.queue:
			resp, err := worker.client.Get(p.Address())
			if err != nil {
				p.Ready() <- true
			} else {
				p.Build(resp)
				worker.send(p.Links())
			}

		case <-worker.quit:
			return
		}
	}
}

func (worker *Worker) Stop() {
	worker.quit <- true
}

func (worker *Worker) Overload() chan bool {
	return worker.overload
}

func (worker *Worker) send(links []job.Job) {
	for i := range links {
		select {
		case worker.record <- links[i]:
		default:
			worker.Overload() <- true
		}
	}
}

func NewWorker(queue <-chan job.Job, record chan<- job.Job, client client.HttpClient) *Worker {
	worker := new(Worker)
	worker.client = client
	worker.queue = queue
	worker.record = record
	worker.quit = make(chan bool, 1)
	worker.overload = make(chan bool, 1)
	return worker
}
