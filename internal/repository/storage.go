package repository

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

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

func (s *Storage) Set(val string) (key string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key, err = generateID(s.data)
	if err != nil {
		return "", err
	}

	s.data[key] = val

	return key, nil
}

func (s *Storage) ExistKey(key string) (bool) {
	_, ok := s.data[key]
	
	return ok
}

func generateID(data map[string]string) (string, error) {
	for i := 0; i < 10; i++ {
		id := generateHash()
		_, ok := data[id]
		if !ok {
			return id, nil
		}
	}

	return "", errors.New("ID generate is impossible")
}

func generateHash() (string) {
	rand.New((rand.NewSource(time.Now().UnixNano())))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}