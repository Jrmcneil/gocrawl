package page

type Sync struct {
	record map[string]bool
	out    chan *Page
	in     chan *Page
	quit   chan bool
}

func (sync *Sync) newRecord(page *Page) bool {
	defer func() { sync.record[page.address] = true }()
	_, found := sync.record[page.address]
	return !found
}

func (sync *Sync) Start() {
	go func() {
		for {
			select {
			case page := <-sync.in:
				if sync.newRecord(page) {
					sync.out <- page
				}

			case <-sync.quit:
				return
			}
		}
	}()
}

func NewSync(in chan *Page, out chan *Page, quit chan bool) *Sync {
	return &Sync{
		in: in,
		out: out,
		quit: quit,
		record: make(map[string]bool),
	}
}
