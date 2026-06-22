package middleware

import (
	"time"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func addLogMiddleware(e *echo.Echo, l *zap.SugaredLogger) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			req := c.Request()
			uri := req.RequestURI
			method := req.Method

			err := next(c)

			resp, _ := echo.UnwrapResponse(c.Response())
			respSize := req.Header.Get(echo.HeaderContentLength)
			duration := time.Since(start)

			l.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
				"respCode", resp.Status,
				"respSize", respSize,
				"err", err,
			)

			return err
		}
	})
}
