package job

type Job interface {
	Address() string
	Links() []Job
	Build(string)
	Ready() chan bool
}
