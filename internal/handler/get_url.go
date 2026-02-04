package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func GetGetURLHandler(data *map[string]string) (func(c echo.Context) error) {
	return func (c echo.Context) error {
		id := c.Param("id")

		val, ok := (*data)[id]
		if !ok || val == "" {
			return c.String(http.StatusNotFound, "URL is not found")
		}

		return c.Redirect(http.StatusTemporaryRedirect, val)
	}
}