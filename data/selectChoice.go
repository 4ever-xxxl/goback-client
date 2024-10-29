package data

import (
	"os"
	"time"
)

type SelectChoice struct {
	FileType   FileType
	SuffixName string
	ModifiedAt time.Time
	Size       int64
}

func CheckSelectChoice(file os.FileInfo, choice SelectChoice) bool {
	if file.IsDir() {
		return true
	}

	if choice.FileType != -1 {
		var filetype FileType
		switch file.Mode() & os.ModeType {
		case os.ModeDir:
			filetype = FOLDER
		case os.ModeSymlink:
			filetype = SOFTLINK
		case os.ModeNamedPipe, os.ModeSocket, os.ModeDevice:
			filetype = PIPE
		default:
			filetype = FILE
		}
		if filetype != choice.FileType {
			return false
		}
	}

	if choice.SuffixName != "" {
		if len(file.Name()) < len(choice.SuffixName) || file.Name()[len(file.Name())-len(choice.SuffixName):] != choice.SuffixName {
			return false
		}
	}

	if !choice.ModifiedAt.IsZero() {
		if file.ModTime().After(choice.ModifiedAt) {
			return false
		}
	}

	if choice.Size != 0 {
		if file.Size() > choice.Size {
			return false
		}
	}

	return true
}
