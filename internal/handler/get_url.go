package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/Vaha95/golang_pet/internal/service/audit"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

func GetURLHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		id := c.Param("id")

		s := strategy.GetStrategy(cfg)
		data, err := s.Get(id)
		if err != nil {
			if errors.Is(err, repository.ErrorShortURLKeyNotFound) {
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

		audit.PushToAudit(cfg.Config, DTO.CrateBaseAuditItemFollow(data.URL, cfg.GetUserId()))
		
		if data.DeletedAt != nil {
			return c.JSON(http.StatusGone, data.URL)			
		}

		return c.Redirect(http.StatusTemporaryRedirect, data.URL)
	}
}
