package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func TestGetUrl(t *testing.T) {
	url := "http://njknonnjkn.com"

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}

	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, url)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
	w := httptest.NewRecorder()
	l, _ := getLogger()
	h := GetURLHandler(stCfg, l)

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
