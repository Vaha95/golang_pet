package repository

import (
	"fmt"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
)

type FileStrategy struct {
	cfg config.Config
}

func (s FileStrategy) Save(short string, url string, extId *string) error {
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

func (s FileStrategy) SaveBatch(data []DTO.BatchItem, GenerateHash func() string) error {
	for i := 0; i < len(data); i++ {
		item := &data[i]

		id := GenerateHash()
		err := s.Save(id, item.URL, nil)
		if err != nil {
			return fmt.Errorf("failed to save the short URL to file: %w", err)
		}

		item.Short = id
	}

	return nil
}

func (s FileStrategy) Get(key string) (string, error) {
	data, err := repository.ReadFileStore(s.cfg)
	if err != nil {
		return "", err
	}
	if data == nil {
		data = map[string]string{}
	}
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", repository.ErrorShortURLKeyNotFound, key)
	}

	return val, nil
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