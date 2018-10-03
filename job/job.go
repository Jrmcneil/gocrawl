package job

type Job interface {
	Address() string
	Links() []Job
	Build(string)
	Close()
	Done() chan bool
	Ready() chan bool
}
