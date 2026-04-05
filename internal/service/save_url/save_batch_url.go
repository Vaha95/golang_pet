package saveurl

import (
	"fmt"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

func SaveBatchURL(cfg config.StorageConfig, data []DTO.BatchItem) error {
	strategy := strategy.GetStrategy(cfg)
	userId := cfg.GetUserId()
	err := strategy.SaveBatch(data, &userId, GenerateHash)
	if err != nil {
		return fmt.Errorf("failed to save the short URL to DB: %w", err)
	}

	return nil
}
