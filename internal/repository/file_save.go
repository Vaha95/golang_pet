package repository

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var mu sync.RWMutex

func getFile(fullPath string) (*os.File, error) {
	file, err := os.OpenFile(fullPath, os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v\n", err)
	}

	return file, err
}

func createNewFile(fullPath string) (*os.File, error) {
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v\n", err)
	}

	return file, err
}

// readStore читает map из файла
func ReadStore() (map[string]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	createFileIfNotExist()

	filePath := getFullPath()
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}

	var store map[string]string
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// writeStore записывает map в файл
func WriteStore(store map[string]string) error {
	mu.Lock()
	defer mu.Unlock()

	createFileIfNotExist()

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	filePath := getFullPath()

	return os.WriteFile(filePath, data, 0644)
}