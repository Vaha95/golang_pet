package repository

import (
	"errors"
	"fmt"

	"github.com/Vaha95/golang_pet/internal/config"
)

var ErrorShortURLKeyAlreadyExists = errors.New("short URL key already exists")
var ErrorShortURLKeyNotFound = errors.New("short URL key not found")
var ErrorURLNotFound = errors.New("short URL not found")
var ErrorURLByUserNotFound = errors.New("short URL by user not found")

func SetURL(cfg config.Config, key string, val string) (err error) {
	data, err := ReadFileStore(cfg)
	if err != nil {
		return err
	}

	_, ok := data[key]
	if ok {
		return fmt.Errorf("%w: %s", ErrorShortURLKeyAlreadyExists, key)
	}
	if data == nil {
		data = map[string]string{}
	}

	data[key] = val
	WriteFileStore(cfg, data)

	return nil
}
