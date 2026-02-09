package repository

import (
	"fmt"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

type ShortURLKeyAlreadyExistsError struct{}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (e ShortURLKeyAlreadyExistsError) Error() string {
	return "short URL key already exists"
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
		err = ShortURLKeyAlreadyExistsError{}

		return fmt.Errorf("%w: %s", err, key)
	}

	s.data[key] = val

	return nil
}

func (s *Storage) ExistKey(key string) bool {
	_, ok := s.data[key]

	return ok
}
