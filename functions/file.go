package functions

import (
	"fmt"
	"goback-client/data"
	"os"
	"os/exec"
)

func GetFileInfo(ffilepath string) (data.File, error) {
	fileInfo, err := os.Lstat(ffilepath)
	if err != nil {
		return data.File{}, err
	}

	var filetype data.FileType
	switch fileInfo.Mode() & os.ModeType {
	case os.ModeDir:
		filetype = data.FOLDER
	case os.ModeSymlink:
		filetype = data.SOFTLINK
	case os.ModeNamedPipe:
		filetype = data.PIPE
	case os.ModeSocket:
		filetype = data.PIPE
	case os.ModeDevice:
		filetype = data.PIPE
	default:
		filetype = data.FILE
	}

	// NOTE: WINDOWS 下文件权限是通过ACL来控制的，无法直接获取文件所有者
	cmd := exec.Command("pwsh", "/C", "(", "Get-Acl", ffilepath, ").Owner")
	owner, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行命令时出错: %v\n", err)
		return data.File{}, err
	}

	file := data.File{
		Type:        filetype,
		Name:        fileInfo.Name(),
		Author:      string(owner),
		Path:        ffilepath,
		Permissions: fileInfo.Mode().String(),
		ModifiedAt:  fileInfo.ModTime(),
		MD5:         "",
		InCloud:     false,
		Remark:      "",
	}
	return file, nil
}
