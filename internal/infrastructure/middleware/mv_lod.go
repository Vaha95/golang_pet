package middleware

import (
	"time"
	"go.uber.org/zap"

	echo "github.com/labstack/echo/v4"
)

func AddMiddlewares(e *echo.Echo, l *zap.SugaredLogger) {
	addLogMiddleware(e, l)
}

func addLogMiddleware(e *echo.Echo, l *zap.SugaredLogger) {
	e.Use(func (next echo.HandlerFunc) echo.HandlerFunc {
		logFn := func(c echo.Context) error {
			start := time.Now()

			req := *c.Request()
			uri := req.RequestURI
			method := req.Method

			err := next(c)

			duration := time.Since(start)
			l.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"err", err,
			)

			return err
		}

		return logFn
	})
}