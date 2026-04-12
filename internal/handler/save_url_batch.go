package handler

import (
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	DTO "github.com/Vaha95/golang_pet/internal/model/DTO"
	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func GetSaveURLBatchHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c echo.Context) error {
	return func(c echo.Context) error {
		type APIReqiest struct {
			ExtId string `json:"correlation_id"`
			URL   string `json:"original_url"`
		}
		var req []APIReqiest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		var data []DTO.BatchItem
		for _, v := range req {
			data = append(data, DTO.BatchItem{ExtId: v.ExtId, URL: v.URL})
		}

		err := saveurl.SaveBatchURL(cfg, data)

		if err != nil {
			l.Errorf("Batch url save error: %w", err)

			return c.NoContent(http.StatusInternalServerError)
		}

		type APIResponse struct {
			ExtId string `json:"correlation_id"`
			Short string `json:"short_url"`
		}

		var resp []APIResponse

		for _, v := range data {
			resp = append(resp, APIResponse{
				ExtId: v.ExtId,
				Short: v.Short,
			})
		}

		return c.JSON(http.StatusCreated, resp)
	}
}
