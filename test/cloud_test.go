package test

import (
	"goback-client/data"
	"goback-client/functions"
	"testing"
)

func TestUploadFiles(t *testing.T) {
	var files = []string{"D:\\TESTDIR\\backup\\dir4Test.backup"}
	// 上传文件
	data.Config.JWTTOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MiwidXNlcm5hbWUiOiJSZWQiLCJleHAiOjE3MzAzNTk1MzAsImlzcyI6IkJsdWUifQ.lZNtD9NTFr0BDs_dabFcCLzJZDJULtSpGdNAYokVXds"
	data.Config.RootDir = "0f57b50c-bbf4-4b2d-ae9f-67d5e725bd36"
	err := functions.UploadFiles(files, data.Config.RootDir)
	if err != nil {
		t.Error(err)
	}
}

func TestDownloadFile(t *testing.T) {
	data.Config.JWTTOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MiwidXNlcm5hbWUiOiJSZWQiLCJleHAiOjE3MzAzNTk1MzAsImlzcyI6IkJsdWUifQ.lZNtD9NTFr0BDs_dabFcCLzJZDJULtSpGdNAYokVXds"
	data.Config.RootDir = "0f57b50c-bbf4-4b2d-ae9f-67d5e725bd36"
	var filename = "dir4Test.backup"
	var saveDir = "D:\\TESTDIR\\backup"
	err := functions.DownloadFile(filename, data.Config.RootDir, saveDir)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteFile(t *testing.T) {
	data.Config.JWTTOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MiwidXNlcm5hbWUiOiJSZWQiLCJleHAiOjE3MzAzNTk1MzAsImlzcyI6IkJsdWUifQ.lZNtD9NTFr0BDs_dabFcCLzJZDJULtSpGdNAYokVXds"
	data.Config.RootDir = "0f57b50c-bbf4-4b2d-ae9f-67d5e725bd36"
	var filename = "dir4Test.backup"
	err := functions.DeleteFile(filename, data.Config.RootDir)
	if err != nil {
		t.Error(err)
	}
}
