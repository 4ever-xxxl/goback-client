package test

import (
	"bytes"
	"fmt"
	"goback-client/data"
	"goback-client/functions"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const rootDir = "D:/TESTDIR"
const testDir = "D:/TESTDIR/dir4Test"

func TestFileInfo(t *testing.T) {
	createTestDir(testDir)
	file, err := functions.GetFileInfo(testDir)
	if err != nil {
		t.Error(err)
	}
	t.Log(file)
}

func TestTarUntar(t *testing.T) {
	var buf bytes.Buffer
	// 获取文件信息
	file1, err := functions.GetFileInfo(testDir)
	if err != nil {
		t.Error(err)
	}

	// 打包文件夹
	err = functions.Tar(&file1, &buf)
	if err != nil {
		t.Error(err)
	}

	// 创建临时目录用于解包
	tempDir := rootDir + "/temp"
	err = os.Mkdir(tempDir, os.ModeDir)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 解包到临时目录
	err = functions.UnTar(&buf, tempDir)
	if err != nil {
		t.Error(err)
	}

	// 比较文件夹
	if err = compareDirs(testDir, tempDir+"/dir4Test"); err != nil {
		t.Error(err)
	}
}

func compareDirs(dir1, dir2 string) error {
	return filepath.Walk(dir1, func(path1 string, info1 os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir1, path1)
		if err != nil {
			return err
		}

		path2 := filepath.Join(dir2, relPath)
		info2, err := os.Stat(path2)
		if err != nil {
			return err
		}

		if info1.IsDir() != info2.IsDir() {
			return fmt.Errorf("mismatch: %s is dir: %v, %s is dir: %v", path1, info1.IsDir(), path2, info2.IsDir())
		}

		if !info1.IsDir() {
			data1, err := os.ReadFile(path1)
			if err != nil {
				return err
			}

			data2, err := os.ReadFile(path2)
			if err != nil {
				return err
			}

			if !bytes.Equal(data1, data2) {
				return fmt.Errorf("file content mismatch: %s and %s", path1, path2)
			}
		}

		return nil
	})
}

func TestCompress(t *testing.T) {
	var originBuf, compressedBuf, deCompressBuf bytes.Buffer
	originBuf.WriteString("hello world")
	err := functions.Compress(&originBuf, &compressedBuf)
	if err != nil {
		t.Error(err)
	}
	err = functions.Decompress(&compressedBuf, &deCompressBuf)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, originBuf.String(), deCompressBuf.String())
}

func TestEncryptDecrypt(t *testing.T) {
	var originBuf, encryptedBuf, decryptedBuf bytes.Buffer
	originBuf.WriteString("hello world")
	err := functions.Encrypt(&originBuf, &encryptedBuf, []byte(data.Config.Key))
	if err != nil {
		t.Error(err)
	}
	err = functions.Decrypt(&encryptedBuf, &decryptedBuf, []byte(data.Config.Key))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, originBuf.String(), decryptedBuf.String())
}

func TestBackupRestore(t *testing.T) {
	var originDir = testDir
	fileInfo, err := functions.GetFileInfo(originDir)
	if err != nil {
		t.Error(err)
	}
	if err := functions.Backup(originDir, []byte(data.Config.Key)); err != nil {
		t.Error(err)
	}
	if err := functions.Restore(&fileInfo, []byte(data.Config.Key)); err != nil {
		t.Error(err)
	}
	if err = compareDirs(originDir, "D:/TESTDIR/restore/dir4Test"); err != nil {
		t.Error(err)
	}
}

func createTestDir(dirPath string) error {
	// 创建一些各种类型的文件
	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(dirPath+"/file1.txt ", []byte("hello world1"), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(dirPath+"/file2.txt", []byte("hello world2"), os.ModePerm); err != nil {
		return err
	}
	if err := os.Mkdir(dirPath+"/dir1", os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(dirPath+"/dir1/file3.txt", []byte("hello world3"), os.ModePerm); err != nil {
		return err
	}
	return nil
}
