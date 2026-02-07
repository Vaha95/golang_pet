package handler

import (
	"net/http"

	"github.com/Vaha95/golang_pet/internal/repository"
	echo "github.com/labstack/echo/v4"
)

func GetURLHandler(storage *repository.Storage) (func(c echo.Context) error) {
	return func (c echo.Context) error {
		id := c.Param("id")

		val, ok := storage.Get(id)
		if !ok || val == "" {
			return c.String(http.StatusNotFound, "URL is not found")
		}

		return c.Redirect(http.StatusTemporaryRedirect, val)
	}
}