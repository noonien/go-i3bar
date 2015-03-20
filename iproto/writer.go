package iproto

import (
	"encoding/json"
	"io"
)

type Block struct {
	Name     string `json:"name,omitempty"`
	Instance string `json:"instance,omitempty"`

	FullText  string `json:"full_text,omitempty"`
	ShortText string `json:"short_text,omitempty"`
	Color     string `json:"color,omitempty"`
	MinWidth  int    `json:"min_width,omitempty"`
	Align     string `json:"align,omitempty"`
	Urgent    bool   `json:"urgent,omitempty"`

	Separator           bool `json:"separator,omitempty"`
	SeparatorBlockWidth int  `json:"separator_block_width,omitempty"`
}

type Writer struct {
	w   io.Writer
	enc *json.Encoder
}

func NewWriter(w io.Writer) Writer {
	return Writer{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

func (w Writer) WriteHeader(clicks bool) error {
	header := struct {
		Version     int  `json:"version"`
		ClickEvents bool `json:"click_events,omitempty"`
	}{1, clicks}

	err := w.enc.Encode(header)
	if err != nil {
		return err
	}

	// Write begining of endless array, and an initial empty array
	_, err = w.w.Write([]byte{'[', '[', ']', '\n'})
	return err
}

func (w Writer) WriteBlocks(blocks []*Block) error {
	_, err := w.w.Write([]byte{','})
	if err != nil {
		return err
	}

	return w.enc.Encode(blocks)
}
