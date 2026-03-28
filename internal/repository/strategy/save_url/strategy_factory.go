package repository

import (
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
)

type SaveURLStrategy interface {
	Save(short string, url string, extId *string) error
	SaveBatch(data []DTO.BatchItem, GenerateHash func() string) error
}

func GetStrategy(cfg config.StorageConfig) SaveURLStrategy {
	if cfg.IsDBAllowed {
		return DBStrategy{(cfg.DBService).GetDB()}
	}

	return FileStrategy{cfg.Config}
}