package main

import (
	"goback-client/data"
	panels "goback-client/panels"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentPanel = "Welcome"

var topWindow fyne.Window

func main() {
	a := app.NewWithID("GoBackup")
	w := a.NewWindow("GoBackup")
	topWindow = w
	logLifecycle(a)
	makeTray(a, w)

	w.Resize(fyne.NewSize(900, 600))
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("")
	intro := widget.NewLabel("")
	intro.Wrapping = fyne.TextWrapWord
	setPanel := func(f panels.Panel) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(f.Title)
			topWindow = child
			child.SetContent(f.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(f.Title)
		intro.SetText(f.Intro)

		if f.Title == "Welcome" {
			title.Hide()
			intro.Hide()
		} else {
			title.Show()
			intro.Show()
		}

		content.Objects = []fyne.CanvasObject{f.View(w)}
		content.Refresh()
	}

	panel := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setPanel, false))
	} else {
		split := container.NewHSplit(makeNav(setPanel, true), panel)
		split.Offset = 0.2
		w.SetContent(split)
	}

	w.ShowAndRun()
	defer func() {
		data.SaveLocalFileList(data.LocalFileListPath)
	}()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeTray(a fyne.App, w fyne.Window) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("主窗口", func() {
			log.Println("主窗口 clicked")
			w.Show()
			w.RequestFocus()
		})
		trayMenu := fyne.NewMenu("系统托盘菜单", h)
		desk.SetSystemTrayMenu(trayMenu)
	}
}

func makeNav(setPanel func(f panels.Panel), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return panels.PanelIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := panels.PanelIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := panels.Panels[uid]
			if !ok {
				fyne.LogError("Missing panel panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := panels.Panels[uid]; ok {
				a.Preferences().SetString(preferenceCurrentPanel, uid)
				setPanel(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentPanel, "Welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
