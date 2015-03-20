package iproto

type Widget interface {
	// Initialize widget
	Init(refresh chan bool, stop chan bool) error

	// Render and return an array of blocks
	Render(force bool) (bool, []*Block, error)

	// Return the names for which this widget responds to click events
	Names() []string

	// Handle a click event
	Click(Click) error
}

type Renderer struct {
	Widgets []Widget
	Err     chan error

	cache   map[Widget][]*Block
	refresh chan bool
	stop    chan bool
}

func NewRenderer(widgets []Widget) (*Renderer, error) {
	r := &Renderer{
		Widgets: widgets,

		cache:   make(map[Widget][]*Block, len(widgets)),
		refresh: make(chan bool),
		stop:    make(chan bool),
	}

	for _, widget := range widgets {
		err := widget.Init(r.refresh, r.stop)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Renderer) Start(w Writer) {
	for {
		select {
		case <-r.stop:
			break
		default:
		}

		err := r.render(w)
		if err != nil {
			r.Err <- err
			break
		}

		<-r.refresh
	}
}

func (r *Renderer) render(w Writer) error {
	var blocks []*Block
	for _, widget := range r.Widgets {
		cache, cached := r.cache[widget]
		changed, rblocks, err := widget.Render(!cached)
		if err != nil {
			return err
		}

		wblocks := cache
		if changed {
			wblocks = rblocks
		}

		blocks = append(blocks, wblocks...)
	}

	return w.WriteBlocks(blocks)
}

func (r *Renderer) Stop() {
	close(r.stop)
	close(r.refresh)
}
