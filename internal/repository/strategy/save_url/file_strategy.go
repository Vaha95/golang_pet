package repository

import (
	"fmt"
	"slices"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
)

type FileStrategy struct {
	cfg config.Config
}

func (s FileStrategy) Save(short string, url string, extId *string, userId *int) error {
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return err
	}

	_, ok := data[short]
	if ok {
		return fmt.Errorf("%w: %s", repository.ErrorShortURLKeyAlreadyExists, short)
	}
	if data == nil {
		data = map[string]string{}
	}

	data[short] = url
	repository.WriteFileStore(s.cfg, data)

	return nil
}

func (s FileStrategy) SaveBatch(data []DTO.BatchItem, userId *int, GenerateHash func() string) error {
	for i := 0; i < len(data); i++ {
		item := &data[i]

		id := GenerateHash()
		err := s.Save(id, item.URL, nil, userId)
		if err != nil {
			return fmt.Errorf("failed to save the short URL to file: %w", err)
		}

		item.Short = id
	}

	return nil
}

func (s FileStrategy) Get(key string) (*DTO.ShortItem, error) {
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]string{}
	}
	val, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", repository.ErrorShortURLKeyNotFound, key)
	}

	return &DTO.ShortItem{URL: val}, nil
}

func (s FileStrategy) GetShortByURL(url string) (string, error) {
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", fmt.Errorf("%w: %s", repository.ErrorURLNotFound, url)
	}

	for short, u := range data {
		if u == url {
			return short, nil
		}
	}

	return "", fmt.Errorf("%w: %s", repository.ErrorURLNotFound, url)
}

func (s FileStrategy) GetByUser(userId *int) ([]DTO.ShortItem, error) {
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return make([]DTO.ShortItem, 0), err
	}
	if data == nil {
		return make([]DTO.ShortItem, 0), fmt.Errorf("%w: %d", repository.ErrorURLByUserNotFound, userId)
	}

	result := make([]DTO.ShortItem, len(data))
	for short, u := range data {
		result = append(result, DTO.ShortItem{Short: short, URL: u})
	}

	if len(result) > 0 {
		return result, nil
	}

	return make([]DTO.ShortItem, 0), fmt.Errorf("%w: %d", repository.ErrorURLByUserNotFound, userId)
}

func (s FileStrategy) DeleteBatch(inp DTO.DeleteBatch) error {
	batch := inp.Shorts
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("%w", repository.ErrorURLNotFound)
	}

	for short := range data {
		if slices.Contains(batch, short) {
			delete(data, short)
		}
	}
	repository.WriteFileStore(s.cfg, data)

	return nil
}
