package handler

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

func GetURLByUserHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c echo.Context) error {
	return func(c echo.Context) error {
		s := strategy.GetStrategy(cfg)
		data, err := s.GetByUser(cfg.GetUserId())
		if (err != nil && errors.Is(err, repository.ErrorURLByUserNotFound)) || len(data) <= 0 {
			return c.NoContent(http.StatusNoContent)
		}

		type APIResponse struct {
			Short string `json:"short_url"`
			URL   string `json:"original_url"`
		}

		var resp []APIResponse
		for _, v := range data {
			short, err := url.JoinPath(cfg.Config.URLHost, v.Short)
			if err != nil {
				l.Errorf("Can`t to create short URL: %w", err)

				continue
			}

			resp = append(resp, APIResponse{
				URL: v.URL,
				Short: short,
			})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
