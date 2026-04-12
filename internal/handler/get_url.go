package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

func GetURLHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		id := c.Param("id")

		s := strategy.GetStrategy(cfg)
		data, err := s.Get(id)
		if err != nil && errors.Is(err, repository.ErrorShortURLKeyNotFound) {
			return c.JSON(http.StatusNotFound, "URL is not found")
		}
		if data.DeletedAt != "" {
			return c.NoContent(http.StatusGone)			
		}
		if data.URL == "" {
			l.Errorf("Url is empty")

			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Redirect(http.StatusTemporaryRedirect, data.URL)
	}
}
