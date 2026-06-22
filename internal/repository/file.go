package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Vaha95/golang_pet/internal/config"
)

// File system paths used by the file-based storage.
const (
	// MainDir is the root directory for local file storage.
	MainDir = `./var`
	// FilenameDirPrefix is the full directory path where data files are stored.
	FilenameDirPrefix = MainDir + `/storage/data`
	// Filename is the default name of the JSON data file.
	Filename = `urls-data.json`
)

func createFileIfNotExist(cfg config.Config) error {
	fullPath := getFullPath(cfg)
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

func getFullPath(cfg config.Config) string {
	filePath := cfg.FilePath
	if filePath == `` {
		filePath = filepath.Join(FilenameDirPrefix, "/", Filename)
	}

	return filePath
}

func createNewFile(fullPath string) (*os.File, error) {
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("can`t storage file: %w", err)
	}

	defValue, err := json.Marshal(map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("can`t storage file: %w", err)
	}
	file.Write(defValue)

	return file, err
}

func getFile(fullPath string) (*os.File, error) {
	file, err := os.OpenFile(fullPath, os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("can`t storage file: %w", err)
	}

	return file, err
}
