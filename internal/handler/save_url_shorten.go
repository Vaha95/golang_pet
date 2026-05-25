package handler

import (
	"errors"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	"github.com/Vaha95/golang_pet/internal/service/audit"
	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func GetSaveURLShortenHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		type APIReqiest struct {
			URI string `json:"url"`
		}
		var data APIReqiest

		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		type APIResponse struct {
			Result string `json:"result"`
		}

		path, err := saveurl.SaveURL(cfg, data.URI)
		if err != nil {
			if errors.Is(err, saveurl.ErrorUrlAlreadyExists) {
				return c.JSON(http.StatusConflict, APIResponse{Result: path})
			} else if errors.Is(err, saveurl.ErrorSaveToStorage) {
				l.Errorf("Shorten url save error: %w", err)

				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		}
		audit.PushToAudit(cfg.Config, DTO.CrateBaseAuditItemShorten(data.URI, *cfg.GetUserId()))

		return c.JSON(http.StatusCreated, APIResponse{Result: path})
	}
}
