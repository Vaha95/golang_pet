package repository

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

const (
	FilenameDirPrefix = `./var/storage/data`;
	Filename = `urls-data.json`
)

func createFileIfNotExist () {
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

	defer file.Close()
}

func getFullPath () string {
	return filepath.Join(FilenameDirPrefix, "/", Filename)
}