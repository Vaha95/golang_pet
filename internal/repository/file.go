package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	MainDir = `./var`
	FilenameDirPrefix = MainDir + `/storage/data`;
	Filename = `urls-data.json`
)

func createFileIfNotExist () error {
	fullPath := getFullPath()
	dir := filepath.Dir(fullPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755) 
		if err != nil {
			log.Fatal(err)
		}
	}

	var file *os.File
	_, err := os.Stat(fullPath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		file, err = createNewFile(fullPath)
	} else {
		file, err = getFile(fullPath)
	}

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func getFullPath () string {
	return filepath.Join(FilenameDirPrefix, "/", Filename)
}

func createNewFile(fullPath string) (*os.File, error) {
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("can`t storage file: %w", err)
	}

	return file, err
}

func getFile(fullPath string) (*os.File, error) {
	file, err := os.OpenFile(fullPath, os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("can`t storage file: %w", err)
	}

	return file, err
}