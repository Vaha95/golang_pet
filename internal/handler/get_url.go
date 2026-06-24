package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

// GetURLHandler returns an Echo handler that resolves a short URL and redirects to the original.
func GetURLHandler(cfg config.StorageConfig, l *zap.SugaredLogger, auditCh chan DTO.BaseAuditItem) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		id := c.Param("id")

		s := strategy.GetStrategy(cfg)
		data, err := s.Get(id)
		if err != nil {
			if errors.Is(err, repository.ErrorShortURLKeyNotFound) || errors.Is(err, strategy.ErrorUrlNotFound) {
				return c.JSON(http.StatusNotFound, "URL is not found")
			}

			l.Errorf("Failed to get url from DB: %w", err)

			return c.NoContent(http.StatusInternalServerError)
		}

		if data == nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		if data.URL == "" {
			l.Errorf("Url is empty")

			return c.NoContent(http.StatusInternalServerError)
		}

		auditCh <- DTO.CreateBaseAuditItemFollow(data.URL, cfg.GetUserId())

		if data.DeletedAt != nil {
			return c.JSON(http.StatusGone, data.URL)
		}

		return c.Redirect(http.StatusTemporaryRedirect, data.URL)
	}
}
