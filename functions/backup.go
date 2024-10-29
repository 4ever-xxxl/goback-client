package functions

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"goback-client/data"
	"log"
	"os"
)

func Backup(ffilepath string, key []byte, opts ...data.SelectChoice) error {
	var options data.SelectChoice
	if len(opts) > 0 {
		options = opts[0]
	}
	for _, file := range data.LocalFileList {
		if file.Path == ffilepath {
			return fmt.Errorf("文件已存在")
		}
	}
	newFileInfo, err := GetFileInfo(ffilepath)
	if err != nil {
		return err
	}
	newFileInfo.SelectChoice = options

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
	log.Println(newFileInfo.Name + " 文件已备份")
	return nil
}

func ReBackup(f data.File, key []byte) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	var tarredBuf bytes.Buffer
	if err = Tar(&f, &tarredBuf); err != nil {
		return err
	}

	hash := md5.Sum(tarredBuf.Bytes())
	if f.MD5 == fmt.Sprintf("%x", hash) {
		log.Println(f.Name + " 文件未发生变化，无需重新备份")
		return nil
	}
	f.MD5 = fmt.Sprintf("%x", hash)

	var compressedBuf bytes.Buffer
	if err = Compress(&tarredBuf, &compressedBuf); err != nil {
		return err
	}

	var encryptedBuf bytes.Buffer
	if err = Encrypt(&compressedBuf, &encryptedBuf, key); err != nil {
		return err
	}

	backupPath := data.Config.BackupDir + f.Name + ".backup"
	if err = os.Remove(backupPath); err != nil {
		return err
	}
	file, err = os.Create(backupPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(encryptedBuf.Bytes()); err != nil {
		return err
	}
	log.Println(f.Name + " 文件已重新备份")
	return nil
}
