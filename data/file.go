package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// FileType 枚举类型表示文件的不同类型
type FileType int

const (
	FILE FileType = iota
	FOLDER
	PIPE
	HARDLINK
	SOFTLINK
)

type File struct {
	Type        FileType
	Name        string
	Author      string
	Path        string
	Permissions string
	ModifiedAt  time.Time
	MD5         string
	Remark      string
	//  DIY Backup
	SelectChoice SelectChoice
	// For Cloud
	InCloud bool
	SHA256  string
} // 保存了文件的基本元数据

type CloudFile struct {
	UUID      string
	Filename  string
	Filesize  int
	Ext       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var LocalFileList = []File{}      // 保存了本地文件列表
var CloudFileList = []CloudFile{} // 保存了云端文件列表
const LocalFileListPath = "static/local_file_list.json"

func (f FileType) String() string {
	switch f {
	case FILE:
		return "文件"
	case FOLDER:
		return "文件夹"
	case PIPE:
		return "管道"
	case HARDLINK:
		return "硬链接"
	case SOFTLINK:
		return "软链接"
	default:
		return "未知类型"
	}
}

func GetFileType(name string) FileType {
	switch name {
	case "":
		return -1
	case "文件":
		return FILE
	case "文件夹":
		return FOLDER
	case "管道":
		return PIPE
	case "硬链接":
		return HARDLINK
	case "软链接":
		return SOFTLINK
	default:
		return -1
	}
}

func SaveLocalFileList(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(LocalFileList); err != nil {
		return fmt.Errorf("无法编码文件列表: %v", err)
	}

	return nil
}

func LoadLocalFileList(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			LocalFileList = []File{}
			return nil
		}
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&LocalFileList); err != nil {
		return fmt.Errorf("无法解码文件列表: %v", err)
	}

	return nil
}

func init() {
	LoadLocalFileList(LocalFileListPath)
}
