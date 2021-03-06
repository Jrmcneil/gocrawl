package record

import (
	"gocrawl/job"
)

type Record struct {
	record map[string]bool
	Out    chan<- job.Job
	In     <-chan job.Job
	quit   chan bool
}

func (record *Record) newRecord(page job.Job) bool {
	defer func() { record.record[page.Address()] = true }()
	_, found := record.record[page.Address()]
	return !found
}

func (record *Record) Start() {
	for {
		select {
		case page := <-record.In:
			if record.newRecord(page) {
				record.Out <- page
			} else {
				page.Ready() <- true
			}

		case <-record.quit:
			return
		}
	}
}

func (record *Record) Stop() {
	record.quit <- true
}

func NewRecord(in <-chan job.Job, out chan<- job.Job, quit chan bool) *Record {
	return &Record{
		In:     in,
		Out:    out,
		quit:   quit,
		record: make(map[string]bool),
	}
}
