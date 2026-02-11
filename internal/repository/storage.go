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

var ShortURLKeyAlreadyExistsError = errors.New("short URL key already exists")

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	val, ok := s.data[key]

	return val, ok
}

func (s *Storage) Set(key string, val string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]
	if ok {
		return fmt.Errorf("%w: %s", ShortURLKeyAlreadyExistsError, key)
	}

	s.data[key] = val

	return nil
}
