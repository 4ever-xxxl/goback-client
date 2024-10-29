package panels

import (
	"goback-client/data"
	"goback-client/functions"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var config = data.Config

func settingsScreen(win fyne.Window) fyne.CanvasObject {
	log.Println(config)
	backupDirBinding := binding.BindString(&config.BackupDir)
	backupDirEntry := widget.NewEntryWithData(backupDirBinding)
	backupDirFormItem := widget.NewFormItem("备份目录", backupDirEntry)

	restoreDirBinding := binding.BindString(&config.RestoreDir)
	restoreDirEntry := widget.NewEntryWithData(restoreDirBinding)
	restoreDirFormItem := widget.NewFormItem("恢复目录", restoreDirEntry)

	keyBinding := binding.BindString(&config.Key)
	keyEntry := widget.NewEntryWithData(keyBinding)
	keyEntry.Password = true
	keyFormItem := widget.NewFormItem("加密密钥", keyEntry)

	cloudBinding := binding.BindString(&config.Cloud)
	cloudEntry := widget.NewEntryWithData(cloudBinding)
	cloudFormItem := widget.NewFormItem("云端地址", cloudEntry)

	portBinding := binding.BindString(&config.Port)
	portEntry := widget.NewEntryWithData(portBinding)
	portFormItem := widget.NewFormItem("端口", portEntry)

	form := &widget.Form{
		Items: []*widget.FormItem{
			backupDirFormItem,
			restoreDirFormItem,
			keyFormItem,
			cloudFormItem,
			portFormItem,
		},
	}

	restoreToOriginalBinding := binding.BindBool(&config.RestoreToOriginal)
	restoreToOriginalCheckbox := widget.NewCheckWithData("恢复到原始目录", restoreToOriginalBinding)

	timeBackupBinding := binding.BindBool(&config.TimedBackup)
	timeBackupCheckbox := widget.NewCheckWithData("定时备份", timeBackupBinding)
	timeBackupCheckbox.OnChanged = func(value bool) {
		config.TimedBackup = value
		if value {
			functions.TimedBackup(1 * time.Hour) // 定时备份间隔为 1 小时
		} else {
			if functions.StopTimedBackup != nil {
				close(functions.StopTimedBackup)
				functions.StopTimedBackup = nil
			}
		}
	}

	fsNotifyBinding := binding.BindBool(&config.FsNotify)
	fsNotifyCheckbox := widget.NewCheckWithData("文件系统感知", fsNotifyBinding)

	saveButton := widget.NewButton("保存设置", func() {
		data.Config = config
		log.Println(data.Config)
	})

	resetButton := widget.NewButton("重置设置", func() {
		config.Init()
		data.Config.Init()
	})

	return container.NewVBox(
		form,
		restoreToOriginalCheckbox,
		timeBackupCheckbox,
		fsNotifyCheckbox,
		saveButton,
		resetButton,
	)
}
