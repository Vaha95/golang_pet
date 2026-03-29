package handler

import (
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
	"github.com/labstack/echo/v4"
)

func GetSaveURLBatchHandler(cfg config.StorageConfig) (func(c echo.Context) error) {
	return func(c echo.Context) error {
		type APIReqiest struct {
			ExtId string `json:"correlation_id"`
			URL string `json:"original_url"`
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
			return c.JSON(http.StatusInternalServerError, err.Error()) 
		}

		return c.JSON(http.StatusCreated, data)
	}
}