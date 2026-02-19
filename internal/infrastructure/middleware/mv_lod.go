package middleware

import (
	"time"
	"go.uber.org/zap"
	"github.com/labstack/echo/v4"
)

func addLogMiddleware(e *echo.Echo, l *zap.SugaredLogger) {
	e.Use(func (next echo.HandlerFunc) echo.HandlerFunc {
		logFn := func(c echo.Context) error {
			start := time.Now()

			req := *c.Request()
			uri := req.RequestURI
			method := req.Method

			err := next(c)

			respCode := c.Response().Status
			respSize := c.Response().Size
			duration := time.Since(start)
			
			l.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"respCode", respCode,
				"respSize", respSize,
				"err", err,
			)

			return err
		}

		return logFn
	})
}