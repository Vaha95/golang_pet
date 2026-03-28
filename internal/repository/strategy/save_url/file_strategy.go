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
		err := s.Save(GenerateHash(), data[i].URL, nil)
		if err != nil {
			return fmt.Errorf("failed to save the short URL to file: %w", err)
		}
	}

	return nil
}