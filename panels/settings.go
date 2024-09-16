package panels

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {
	log.Println("TODO: settingsScreen")
	return container.NewCenter(widget.NewLabel("Settings"))
}
