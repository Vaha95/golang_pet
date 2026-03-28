package repository

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/Vaha95/golang_pet/internal/config"
)

var rmu sync.RWMutex

func ReadFileStore(cfg config.Config) (map[string]string, error) {
	rmu.RLock()
	defer rmu.RUnlock()

	err := createFileIfNotExist(cfg)
	if err != nil {
		return nil, err
	}

	filePath := getFullPath(cfg)
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