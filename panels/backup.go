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

var selectedID int = -1

func backupScreen(win fyne.Window) fyne.CanvasObject {
	fdbox := fileDetaileHbox()
	list := widget.NewList(
		func() int {
			return len(data.LocalFileList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(getFileIcon(data.LocalFileList[id].Type))
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data.LocalFileList[id].Name)
		},
	)

	addFolderButton := widget.NewButton("AddFolder", func() {
		dialog.ShowFolderOpen(func(listURI fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if listURI == nil {
				return
			}
			// 请输入加密密钥
			passwordEntry := widget.NewPasswordEntry()
			form := &widget.Form{
				Items: []*widget.FormItem{
					widget.NewFormItem("密钥", passwordEntry),
				},
			}
			dialog.ShowForm("加密密钥", "确认", "取消", form.Items, func(b bool) {
				if !b {
					return
				}
				log.Println("Password: ", passwordEntry.Text)
				if err = functions.Backup(listURI.Path(), []byte(passwordEntry.Text)); err != nil {
					dialog.ShowError(err, win)
					return
				}
				list.Refresh()
			}, win)
		}, win)
	})
	addFileButton := widget.NewButton("AddFile", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			if err = functions.Backup(reader.URI().Path(), data.Key); err != nil {
				dialog.ShowError(err, win)
				return
			}
			list.Refresh()
		}, win)
		fd.Show()
	})
	restoreButton := widget.NewButton("Restore", func() {
		if selectedID == -1 {
			dialog.ShowInformation("No File Selected", "Please select a file to restore", win)
			return
		}
		if err := functions.Restore(&data.LocalFileList[selectedID], data.Key); err != nil {
			dialog.ShowError(err, win)
			return
		}
		dialog.ShowInformation("Restore", "Restore completed", win)
	})
	deleteButton := widget.NewButton("Delete", func() {
		dialog.ShowConfirm("Delete", "Are you sure you want to delete the selected files?", func(b bool) {
			if b {
				// TODO: 删除按钮的点击事件处理
				log.Println("Delete confirmed")
				data.LocalFileList = append(data.LocalFileList[:selectedID], data.LocalFileList[selectedID+1:]...)
				selectedID = -1
				fileDetaileRefresh(fdbox, nil)
			}
		}, win)
	})
	unselectButton := widget.NewButton("Unselect", func() {
		list.UnselectAll()
	})

	restoreButton.Disable()
	deleteButton.Disable()
	unselectButton.Disable()

	list.OnSelected = func(id widget.ListItemID) {
		list.Select(id)
		selectedID = id
		fileDetaileRefresh(fdbox, &data.LocalFileList[id])
		restoreButton.Enable()
		deleteButton.Enable()
		unselectButton.Enable()
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedID = -1
		fileDetaileRefresh(fdbox, nil)
		restoreButton.Disable()
		deleteButton.Disable()
		unselectButton.Disable()
	}

	buttons := container.NewHBox(addFolderButton, addFileButton, restoreButton, deleteButton, unselectButton)
	leftPanel := container.NewVSplit(buttons, list)
	leftPanel.Offset = 0.1
	content := container.NewHSplit(leftPanel, container.NewCenter(fdbox))
	content.Offset = 0.7
	return content
}

func fileDetaileHbox() fyne.CanvasObject {
	widgets := make([]fyne.CanvasObject, 11)
	widgets[0] = widget.NewIcon(theme.DocumentIcon())
	for i := 1; i < 11; i++ {
		widgets[i] = widget.NewLabel("")
		widgets[i].(*widget.Label).Alignment = fyne.TextAlignLeading
	}
	widgets[1].(*widget.Label).SetText("Select An Item From The List")
	return container.NewVBox(widgets...)
}

func fileDetaileRefresh(fdbox fyne.CanvasObject, f *data.File) {
	vbox := fdbox.(*fyne.Container)
	labels := vbox.Objects
	if f == nil {
		labels[0].(*widget.Icon).SetResource(theme.DocumentIcon())
		labels[1].(*widget.Label).SetText("Select An Item From The List")
		for i := 2; i < len(labels); i++ {
			labels[i].(*widget.Label).SetText("")
		}
	} else {
		labels[0].(*widget.Icon).SetResource(getFileIcon(f.Type))
		labels[1].(*widget.Label).SetText("文件名: " + f.Name)
		labels[2].(*widget.Label).SetText("文件类型: " + f.Type.String())
		labels[3].(*widget.Label).SetText("作者: " + f.Author)
		labels[4].(*widget.Label).SetText("路径: " + f.Path)
		labels[5].(*widget.Label).SetText("权限: " + f.Permissions)
		labels[7].(*widget.Label).SetText("修改时间: " + f.ModifiedAt.Format("2006-01-02 15:04:05"))
		labels[8].(*widget.Label).SetText("MD5: " + f.MD5)
		labels[9].(*widget.Label).SetText("是否在云端: " + strconv.FormatBool(f.InCloud))
		labels[10].(*widget.Label).SetText("备注: " + f.Remark)
	}
}

func getFileIcon(fileType data.FileType) fyne.Resource {
	switch fileType {
	case data.FILE:
		return theme.DocumentIcon()
	case data.FOLDER:
		return theme.FolderIcon()
	case data.PIPE:
		return data.PipeLogo
	case data.HARDLINK:
		return data.HardLinkLogo
	case data.SOFTLINK:
		return data.SoftLinkLogo
	default:
		return theme.DocumentIcon()
	}
}
