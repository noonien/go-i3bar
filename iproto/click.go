package iproto

import (
	"bufio"
	"encoding/json"
	"io"
)

type Click struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`

	Button int `json:"button"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type ClickHandler struct {
	Widgets []Widget
	Err     chan error

	handlers map[string][]Widget
	stop     bool
}

func (ch *ClickHandler) Listen(r io.Reader) {
	ch.init()

	buf := bufio.NewReader(r)
	dec := json.NewDecoder(buf)
	for {
		err := discardPreamble(buf)
		if err != nil {
			ch.Err <- err
			break
		}

		var click Click
		err = dec.Decode(&click)
		if ch.stop {
			break
		}
		if err != nil {
			ch.Err <- err
			break
		}

		err = ch.handleClick(click)
		if err != nil {
			ch.Err <- err
			break
		}
	}
}

func (ch *ClickHandler) init() {
	handlers := make(map[string][]Widget)

	for _, widget := range ch.Widgets {
		names := widget.Names()
		for _, name := range names {
			widgets, _ := handlers[name]
			handlers[name] = append(widgets, widget)
		}
	}

	ch.Err = make(chan error)
	ch.handlers = handlers
	ch.stop = false
}

func (ch *ClickHandler) handleClick(c Click) error {
	widgets, _ := ch.handlers[c.Name]
	for _, widget := range widgets {
		err := widget.Click(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func discardPreamble(buf *bufio.Reader) error {
	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			return err
		}

		if r == '{' {
			buf.UnreadRune()
			return nil
		}
	}
}

func (ch *ClickHandler) Stop() {
	ch.stop = true
}
