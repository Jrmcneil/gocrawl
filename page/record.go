package page

type Record struct {
	record map[string]bool
	out    chan *Page
	in     chan *Page
	quit   chan bool
}

func (record *Record) newRecord(page *Page) bool {
	defer func() { record.record[page.address] = true }()
	_, found := record.record[page.address]
	return !found
}

func (record *Record) Start() {
	go func() {
		for {
			select {
			case page := <-record.in:
				if record.newRecord(page) {
					record.out <- page
				} else {
					page.done <- true
				}

			case <-record.quit:
				return
			}
		}
	}()
}

func NewRecord(in chan *Page, out chan *Page, quit chan bool) *Record {
	return &Record{
		in:     in,
		out:    out,
		quit:   quit,
		record: make(map[string]bool),
	}
}
