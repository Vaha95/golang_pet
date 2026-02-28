package middleware

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var ErrorZapLoggerInitialize = errors.New("init logger is failed")

func AddMiddlewares(e *echo.Echo) error {
	addEncodeMiddleware(e)

	l, err := getLogger()
	if err != nil {
		return err
	}
	addLogMiddleware(e, l)

	return nil
}

func getLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("%w", ErrorZapLoggerInitialize)
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}