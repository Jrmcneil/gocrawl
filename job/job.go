package job

type Job interface {
	Address() string
	Links() []Job
	ResetLinks()
	Build(string)
	Ready() chan bool
}
