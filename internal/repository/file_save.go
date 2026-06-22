package repository

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Vaha95/golang_pet/internal/config"
)

var mu sync.RWMutex

func WriteFileStore(cfg config.Config, store map[string]string) error {
	mu.Lock()
	defer mu.Unlock()

	err := createFileIfNotExist(cfg)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	filePath := getFullPath(cfg)

	return os.WriteFile(filePath, data, 0644)
}
