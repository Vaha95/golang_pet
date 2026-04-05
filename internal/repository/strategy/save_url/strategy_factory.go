package repository

import (
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
)

type URLStrategy interface {
	Save(short string, url string, extId *string, userId *int) error
	SaveBatch(data []DTO.BatchItem, userId *int, GenerateHash func() string) error
	Get(key string) (string, error)
	GetShortByURL(url string) (string, error)
	GetByUser(userId *int) ([]DTO.ShortItem, error)
}

func GetStrategy(cfg config.StorageConfig) URLStrategy {
	if cfg.IsDBAllowed {
		return DBStrategy{(cfg.DBService).GetDB(), cfg.Config}
	}

	return FileStrategy{cfg.Config}
}