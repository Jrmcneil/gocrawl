package worker

type Worker struct {
	client  HttpClient
	queue   chan Job
	record  chan Job
	quit    chan bool
}

type HttpClient interface {
	Get(address string) (string, error)
}

type Job interface {
	Address() string
	Links() []Job
	Build(string)
	Close()
}

func (worker *Worker) Start() {
	go func() {
		for {
			select {
			case p := <-worker.queue:
				resp, err := worker.client.Get(p.Address())
				if err != nil {
					p.Close()
				} else {
					p.Build(resp)
					worker.send(p.Links())
				}

			case <-worker.quit:
				return
			}
		}
	}()
}

func (worker *Worker) send(links []Job) {
	for _, link := range links {
		worker.record <- link
	}
}

func NewWorker(queue chan Job, record chan Job, client HttpClient) *Worker {
	worker := new(Worker)
	worker.client = client
	worker.queue = queue
	worker.record = record
	worker.quit = make(chan bool, 1)
	return worker
}
