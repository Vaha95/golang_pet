package repository

import (
	"errors"
	"fmt"

	"github.com/Vaha95/golang_pet/internal/config"
)

// ErrorShortURLKeyAlreadyExists is returned when the generated short key is already in use.
var ErrorShortURLKeyAlreadyExists = errors.New("short URL key already exists")

// ErrorShortURLKeyNotFound is returned when looking up a short key that doesn't exist.
var ErrorShortURLKeyNotFound = errors.New("short URL key not found")

// ErrorURLNotFound is returned when the target URL is not in storage.
var ErrorURLNotFound = errors.New("short URL not found")

// ErrorURLByUserNotFound is returned when the user has no saved short URLs.
var ErrorURLByUserNotFound = errors.New("short URL by user not found")

// SetURL adds a short-to-URL mapping to the file store.
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
