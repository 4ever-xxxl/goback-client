package panels

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func aboutScreen(_ fyne.Window) fyne.CanvasObject {
	log.Println("TODO: aboutScreen")
	return container.NewCenter(widget.NewLabel("About"))
}
