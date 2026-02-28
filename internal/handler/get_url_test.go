package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/Vaha95/golang_pet/internal/service/save_url"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrl(t *testing.T) {
	url := "http://vfdfbdfbd.com"

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost: `http://localhost:8080`,
		FilePath: ``,
	}

	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, url)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
	w := httptest.NewRecorder()
	h := GetURLHandler(cfg)

	c := echo.New().NewContext(request, w)
	c.SetParamNames("id")
	c.SetParamValues(id)

	h(c)

	res := w.Result()
	assert.Equal(t, 307, res.StatusCode)

	defer res.Body.Close()

	loc := res.Header.Get("Location")
	require.Equal(t, url, loc)
}