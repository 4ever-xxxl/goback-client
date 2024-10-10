package functions

import (
	"archive/tar"
	"bytes"
	"goback-client/data"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Tar(f *data.File, buf *bytes.Buffer) error {
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		log.Println(header)

		header.Name, err = filepath.Rel(filepath.Dir(f.Path), path)
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func UnTar(buf *bytes.Buffer, dir string) error {
	tr := tar.NewReader(buf)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		path := filepath.Join(dir, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			if err := os.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
				return err
			}
			continue
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, tr); err != nil {
			return err
		}

		if err := os.Chmod(path, info.Mode()); err != nil {
			return err
		}

		if err := os.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
			return err
		}

	}

	return nil
}
