package deleteurl

import (
	"strings"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

// GetDeleteUrlListener returns a function that runs an async deletion worker.
func GetDeleteUrlListener(cfg config.StorageConfig, deleteCh chan DTO.DeleteBatch, l *zap.SugaredLogger) func() {
	return func() {
		for {
			msg := <-deleteCh
			deleteURLs(cfg, msg, l)
		}
	}
}

func deleteURLs(cfg config.StorageConfig, data DTO.DeleteBatch, l *zap.SugaredLogger) {
	l.Infof("Delete process start. Short %s", strings.Join(data.Shorts, ", "))
	strategy := strategy.GetStrategy(cfg)
	err := strategy.DeleteBatch(data)
	if err != nil {
		l.Error(err)
	}
}
