package repository

import (
	"encoding/json"
	"os"
	"sync"
)

var mu sync.RWMutex

func WriteStore(store map[string]string) error {
	mu.Lock()
	defer mu.Unlock()

	err := createFileIfNotExist()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	filePath := getFullPath()

	return os.WriteFile(filePath, data, 0644)
}