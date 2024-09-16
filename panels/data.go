package panels

import "fyne.io/fyne/v2"

type Panel struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// panels contains all the panels that can be displayed
	Panels = map[string]Panel{
		"Welcome":  {"Welcome", "", welcomeScreen},
		"Backup":   {"Backup", "请选择你需要备份/恢复的文件", backupScreen},
		"Cloud":    {"Cloud", "", cloudScreen},
		"Settings": {"Settings", "", settingsScreen},
		"About":    {"About", "", aboutScreen},
	}

	// CurrentPanelIndex is a map of panel names to their index in the panels list
	PanelIndex = map[string][]string{
		"": {"Welcome", "Backup", "Cloud", "Settings", "About"},
	}
)
