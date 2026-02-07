package repository

import "sync"

type Storage struct {
	mu sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	val, ok := s.data[key]

	return val, ok
}

func (s *Storage) Set(key string, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
}

func (s *Storage) ExistKey(key string) (bool) {
	_, ok := s.data[key]
	
	return ok
}