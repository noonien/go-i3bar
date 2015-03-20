package main

import (
	"log"
	"os"

	"github.com/noonien/i3bar"
	"github.com/noonien/i3bar/widgets"
)

func main() {
	writer := bar.NewWriter(os.Stdout)

	err := writer.WriteHeader(true)
	if err != nil {
		log.Fatal(err)
	}

	widgets := []bar.Widget{
		&widgets.TimeWidget{
			Format: "15:04:05",
		},
	}

	renderer, err := bar.NewRenderer(widgets)
	if err != nil {
		log.Fatal(err)
	}
	go renderer.Start(writer)

	chandler := bar.ClickHandler{Widgets: widgets}
	go chandler.Listen(os.Stdin)

	for {
		select {
		case err := <-renderer.Err:
			_ = err

		case err := <-chandler.Err:
			_ = err
		}
	}
}
