package functions

import (
	"goback-client/data"
	"os"
	"path/filepath"
)

func Delete(f data.File) error {
	backupFilePath := filepath.Join(data.Config.BackupDir, f.Name)
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return err
	}
	err := os.Remove(backupFilePath)
	if err != nil {
		return err
	}
	return nil
}
