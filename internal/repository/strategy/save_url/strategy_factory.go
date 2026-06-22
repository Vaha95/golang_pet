package repository

import (
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
)

// URLStrategy defines the operations for storing and retrieving short URLs.
type URLStrategy interface {
	Save(short string, url string, extId *string, userId *int) error
	SaveBatch(data []DTO.BatchItem, userId *int, GenerateHash func() string) error
	Get(key string) (*DTO.ShortItem, error)
	GetShortByURL(url string) (string, error)
	GetByUser(userId *int) ([]DTO.ShortItem, error)
	DeleteBatch(DTO.DeleteBatch) error
}

// GetStrategy returns the URLStrategy implementation based on whether DB is available.
func GetStrategy(cfg config.StorageConfig) URLStrategy {
	if cfg.IsDBAllowed {
		return DBStrategy{(cfg.DBService).GetDB(), cfg.Config}
	}

	return FileStrategy{cfg.Config}
}
