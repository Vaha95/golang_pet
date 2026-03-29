package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type DBService interface {
    Ping() error
}

func GetPingDBHandler(service DBService) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		err := service.Ping()
		if err != nil {
			log.Errorf("Fail ping to DB: %v", err)

			return c.JSON(http.StatusInternalServerError, nil)
		}

		return c.JSON(http.StatusOK, nil)
	}
}