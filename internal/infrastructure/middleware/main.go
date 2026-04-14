package middleware

import (
	"errors"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

var ErrorZapLoggerInitialize = errors.New("init logger is failed")

func AddMiddlewares(cfg config.StorageConfig, e *echo.Echo, l *zap.SugaredLogger) error {
	addAuthMiddleware(cfg, e, l)
	addEncodeMiddleware(e)
	addLogMiddleware(e, l)

	return nil
}
