package functions

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"goback-client/data"
	"os"
)

func Backup(ffilepath string, key []byte) error {
	for _, file := range data.LocalFileList {
		if file.Path == ffilepath {
			return fmt.Errorf("文件已存在")
		}
	}
	newFileInfo, err := GetFileInfo(ffilepath)
	if err != nil {
		return err
	}

	var tarredBuf bytes.Buffer
	if err = Tar(&newFileInfo, &tarredBuf); err != nil {
		return err
	}

	hash := md5.Sum(tarredBuf.Bytes())
	newFileInfo.MD5 = fmt.Sprintf("%x", hash)

	var compressedBuf bytes.Buffer
	if err = Compress(&tarredBuf, &compressedBuf); err != nil {
		return err
	}

	var encryptedBuf bytes.Buffer
	if err = Encrypt(&compressedBuf, &encryptedBuf, key); err != nil {
		return err
	}

	backupPath := data.Config.BackupDir + newFileInfo.Name + ".backup"
	file, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(encryptedBuf.Bytes()); err != nil {
		return err
	}

	data.LocalFileList = append(data.LocalFileList, newFileInfo)

	return nil
}
