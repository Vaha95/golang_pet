package middleware

import (
	"github.com/labstack/echo/v5/middleware"
	"github.com/labstack/echo/v5"
)

func addEncodeMiddleware(e *echo.Echo) {
	e.Use(middleware.Decompress())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}
