package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	echo "github.com/labstack/echo/v4"
)

func TestGetUrl(t *testing.T) {
	url := "http://vfdfbdfbd.com"
	storage := repository.NewStorage()
	id, _ := storage.Set(url)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
	w := httptest.NewRecorder()

	h := GetURLHandler(storage)

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