package panels

import (
	"goback-client/data"
	"goback-client/functions"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var loginFlag = false
var selectedCloudID int = -1

func cloudScreen(w fyne.Window) fyne.CanvasObject {
	contentContainer := container.NewVBox()
	if loginFlag {
		cloudFileList(w, contentContainer)
	} else {
		loginScreen(w, contentContainer)
	}
	return contentContainer
}

func cloudFileList(w fyne.Window, contentContainer *fyne.Container) {
	functions.GetRootDir()
	functions.GetFileList(data.Config.RootDir)

	fdbox := fileDetaileHbox(5)
	list := widget.NewList(
		func() int {
			return len(data.CloudFileList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.DocumentIcon())
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data.CloudFileList[id].Filename)
		},
	)

	refreshButton := widget.NewButton("Refresh", func() {
		functions.GetFileList(data.Config.RootDir)
		list.Refresh()
	})

	deleteButton := widget.NewButton("Delete", func() {
		if selectedCloudID == -1 {
			log.Println("No file selected")
			return
		}
		dialog.ShowConfirm("Delete", "Are you sure you want to delete the selected files?", func(b bool) {
			if b {
				functions.DeleteFile(data.CloudFileList[selectedCloudID].Filename, data.Config.RootDir)
				functions.GetFileList(data.Config.RootDir)
				list.Refresh()
				cloudFileDetaileRefresh(fdbox, nil)
			}
		}, w)
	})

	downloadButton := widget.NewButton("Download", func() {
		if selectedCloudID == -1 {
			log.Println("No file selected")
			return
		}
		if err := functions.DownloadFile(data.CloudFileList[selectedCloudID].Filename, data.Config.RootDir, data.Config.BackupDir); err != nil {
			dialog.ShowError(err, w)
		}
		dialog.ShowInformation("Download", "Download completed", w)
	})

	unselectButton := widget.NewButton("Unselect", func() {
		list.UnselectAll()
	})

	quitButton := widget.NewButton("Quit", func() {
		loginFlag = false
		loginScreen(w, contentContainer)
		contentContainer.Refresh()
	})

	list.OnSelected = func(id widget.ListItemID) {
		list.Select(id)
		selectedCloudID = id
		cloudFileDetaileRefresh(fdbox, &data.CloudFileList[id])
		deleteButton.Enable()
		downloadButton.Enable()
		unselectButton.Enable()
	}

	list.OnUnselected = func(id widget.ListItemID) {
		list.Unselect(id)
		selectedCloudID = -1
		cloudFileDetaileRefresh(fdbox, nil)
		deleteButton.Disable()
		downloadButton.Disable()
		unselectButton.Disable()
	}

	downloadButton.Disable()
	deleteButton.Disable()
	unselectButton.Disable()
	buttons := container.NewHBox(refreshButton, downloadButton, deleteButton, unselectButton, quitButton)
	leftPanel := container.NewVSplit(buttons, list)
	leftPanel.Offset = 0.1
	content := container.NewHSplit(leftPanel, container.NewCenter(fdbox))
	content.Offset = 0.7
	contentContainer.Objects = []fyne.CanvasObject{content}
}

func loginScreen(w fyne.Window, contentContainer *fyne.Container) {
	// 登录界面
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	username.Text = data.Config.Username
	password.Text = data.Config.Password

	login := widget.NewButton("Login", func() {
		data.Config.Username = username.Text
		data.Config.Password = password.Text
		if functions.Login() {
			loginFlag = true
			cloudFileList(w, contentContainer)
			contentContainer.Refresh()
		} else {
			log.Println("Login failed")
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: username},
			{Text: "Password", Widget: password},
		},
	}
	vbox := container.NewVBox(form, login)
	contentContainer.Objects = []fyne.CanvasObject{vbox}
}

func cloudFileDetaileRefresh(fdbox fyne.CanvasObject, f *data.CloudFile) {
	vbox := fdbox.(*fyne.Container)
	labels := vbox.Objects
	if f == nil {
		labels[0].(*widget.Icon).SetResource(theme.DocumentIcon())
		labels[1].(*widget.Label).SetText("Select An Item From The List")
		for i := 2; i < len(labels); i++ {
			labels[i].(*widget.Label).SetText("")
		}
	} else {
		labels[0].(*widget.Icon).SetResource(theme.DocumentIcon())
		labels[1].(*widget.Label).SetText("UUID: " + f.UUID)
		labels[2].(*widget.Label).SetText("文件名: " + f.Filename)
		labels[3].(*widget.Label).SetText("文件大小： " + strconv.Itoa(f.Filesize))
		labels[4].(*widget.Label).SetText("上传时间: " + f.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
