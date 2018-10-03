package page

type Record struct {
	record map[string]bool
	out    chan *Page
	In     chan *Page
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
			case page := <-record.In:
				if record.newRecord(page) {
					record.out <- page
				} else {
					page.Done <- true
				}

			case <-record.quit:
				return
			}
		}
	}()
}

func NewRecord(in chan *Page, out chan *Page, quit chan bool) *Record {
	return &Record{
		In:     in,
		out:    out,
		quit:   quit,
		record: make(map[string]bool),
	}
}
