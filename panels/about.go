package panels

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func aboutScreen(_ fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("GoBack", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	version := widget.NewLabel("Version: 1.0.0")
	author := widget.NewLabel("Author: lx10ng, Blue")
	email := widget.NewLabel("Email: lx10ng@qq.com")
	github := widget.NewLabel("Github: 4ever-xxxl")
	description := widget.NewLabel("Description: A simple backup tool")

	content := container.NewVBox(
		title,
		version,
		author,
		email,
		github,
		description,
	)

	return container.NewCenter(content)
}
