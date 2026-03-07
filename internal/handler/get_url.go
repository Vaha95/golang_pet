package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v4"
)

func GetURLHandler(cfg config.Config) (func(c echo.Context) error) {
	return func (c echo.Context) error {
		id := c.Param("id")

		val, err := repository.GetURLByKey(cfg, id)
		if err != nil && errors.Is(err, repository.ErrorShortURLKeyNotFound) {
			return c.JSON(http.StatusNotFound, "URL is not found")			
		}
		if val == "" {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.Redirect(http.StatusTemporaryRedirect, val)
	}
}