package middleware

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var ErrorZapLoggerInitialize = errors.New("init logger is failed")

func AddMiddlewares(e *echo.Echo, l *zap.SugaredLogger) error {
	addEncodeMiddleware(e)

	addLogMiddleware(e, l)

	return nil
}