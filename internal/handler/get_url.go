package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/repository"
	echo "github.com/labstack/echo/v4"
)

func GetURLHandler() (func(c echo.Context) error) {
	return func (c echo.Context) error {
		id := c.Param("id")

		val, err := repository.GetUrlByKey(id)
		if err != nil && errors.Is(err, repository.ErrorShortURLKeyNotFound) {
			return c.String(http.StatusNotFound, "URL is not found")			
		}
		if val == "" {
			return c.String(http.StatusInternalServerError, err.Error())	
		}

		return c.Redirect(http.StatusTemporaryRedirect, val)
	}
}