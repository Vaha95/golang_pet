package repository

import (
	"encoding/json"
	"os"
	"sync"
)

var rmu sync.RWMutex

func ReadStore() (map[string]string, error) {
	rmu.RLock()
	defer rmu.RUnlock()

	err := createFileIfNotExist()
	if err != nil {
		return nil, err
	}

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