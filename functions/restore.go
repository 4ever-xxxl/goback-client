package functions

import (
	"bytes"
	"goback-client/data"
	"io"
	"os"
)

func Restore(f *data.File, key []byte) error {
	backupPath := data.BackupDir + f.Name + ".backup"

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return err
	}

	var encryptedBuf bytes.Buffer
	encryptedFile, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(&encryptedBuf, encryptedFile); err != nil {
		return err
	}

	var compressedBuf bytes.Buffer
	if err = Decrypt(&encryptedBuf, &compressedBuf, key); err != nil {
		return err
	}

	var tarredBuf bytes.Buffer
	if err = Decompress(&compressedBuf, &tarredBuf); err != nil {
		return err
	}

	var restoreDir string
	if data.RestoreToOriginal {
		restoreDir = f.Path[:len(f.Path)-len(f.Name)]
	} else {
		restoreDir = data.RestoreDir
	}
	if err = UnTar(&tarredBuf, restoreDir); err != nil {
		return err
	}
	return nil
}
