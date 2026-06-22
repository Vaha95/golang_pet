package deleteurl

import (
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

func GetDeleteUrlListener(cfg config.StorageConfig, deleteCh chan DTO.DeleteBatch, l *zap.SugaredLogger) (func()) {
	return func() {
		for {
			msg := <-deleteCh
			deleteURLs(cfg, msg, l)
		}
	}
}

func deleteURLs(cfg config.StorageConfig, data DTO.DeleteBatch, l *zap.SugaredLogger) {
	strategy := strategy.GetStrategy(cfg)
	err := strategy.DeleteBatch(data)
	if err != nil {
		l.Error(err)
	}
}