package fileutil

import (
	"io"
	"os"
	"path/filepath"
)

const path = "Download"

func SaveLocalFile(reader io.Reader, fileName string) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	downloadDir := filepath.Join(rootDir, path)
	_, err = os.Stat(downloadDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, os.ModeAppend)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(filepath.Join(downloadDir, fileName))
	if err != nil {
		return nil
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}
