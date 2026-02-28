package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

var ErrorShortURLKeyAlreadyExists = errors.New("short URL key already exists")
var ErrorShortURLKeyNotFound = errors.New("short URL key already exists")

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func GetUrlByKey(key string) (string, error) {
	data, err := ReadStore()
	if err != nil {
		return "", err
	}
	if data == nil {
		data = map[string]string{}
	}
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrorShortURLKeyNotFound, key)
	}

	return val, nil
}

func SetUrl(key string, val string) (err error) {
	data, err := ReadStore()
	_, ok := data[key]
	if ok {
		return fmt.Errorf("%w: %s", ErrorShortURLKeyAlreadyExists, key)
	}
	if data == nil {
		data = map[string]string{}
	}

	data[key] = val
	WriteStore(data)

	return nil
}
