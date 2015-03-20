package widgets

import (
	"strings"
	"time"

	"github.com/noonien/i3bar"
)

const DefaultTimeFormat = "02-01-2006 15:04"

type TimeWidget struct {
	Format   string
	Location *time.Location
}

func (w *TimeWidget) Init(refresh, stop chan bool) error {
	go func() {
		for {
			select {
			case <-w.tickAtNextUnit():
				refresh <- true
			case <-stop:
				break
			}
		}
	}()

	return nil
}

func (w *TimeWidget) tickAtNextUnit() <-chan time.Time {
	format := w.Format
	if format == "" {
		format = DefaultTimeFormat
	}

	now := time.Now()
	loc := w.Location
	if loc == nil {
		loc = now.Location()
	}

	var next time.Time
	switch true {
	// Seconds
	case strings.Count(format, "5") == 2 || strings.Contains(format, "5"):
		next = now.Add(time.Second)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, loc)

	// Minutes
	case strings.Contains(format, "4"):
		next = now.Add(time.Minute)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, loc)

	// Hours
	case strings.Contains(format, "3") || strings.Contains(format, "15"):
		next = now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, loc)

	// Everything else should be updated daily
	default:
		next = now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, loc)
	}

	return time.After(next.Sub(now))
}

func (w *TimeWidget) Render(force bool) (bool, []*bar.Block, error) {
	format := w.Format
	if format == "" {
		format = DefaultTimeFormat
	}

	return true, []*bar.Block{
		{
			Name:     "time",
			FullText: time.Now().Format(format),
		},
	}, nil
}

func (w *TimeWidget) Names() []string {
	return []string{"time"}
}

func (w *TimeWidget) Click(click bar.Click) error {
	return nil
}
