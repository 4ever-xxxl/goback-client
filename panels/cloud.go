package panels

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func cloudScreen(_ fyne.Window) fyne.CanvasObject {
	log.Println("TODO: cloudScreen")
	return container.NewCenter(widget.NewLabel("Cloud"))
}
