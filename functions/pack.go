package functions

import (
	"bytes"
	"goback-client/data"
	"io/fs"
	"path/filepath"
)

func Pack(f *data.File, buf bytes.Buffer) error {
	if err := filepath.Walk(f.Path, func(path string, fi fs.FileInfo, err error) error {
		return nil
	}); err != nil {
		return err
	}
	return nil
}
