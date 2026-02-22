package middleware

import (
	"go.uber.org/zap"
	"github.com/labstack/echo/v4"
)

func AddMiddlewares(e *echo.Echo, l *zap.SugaredLogger) {
	addEncodeMiddleware(e)
	addLogMiddleware(e, l)
}