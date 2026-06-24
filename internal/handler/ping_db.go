package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// DBService is the minimal database interface required by the ping handler.
type DBService interface {
	Ping() error
}

// GetPingDBHandler returns an Echo handler that checks database connectivity.
func GetPingDBHandler(service DBService, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		err := service.Ping()
		if err != nil {
			l.Errorf("Fail ping to DB: %v", err)

			return c.JSON(http.StatusInternalServerError, nil)
		}

		return c.JSON(http.StatusOK, nil)
	}
}
