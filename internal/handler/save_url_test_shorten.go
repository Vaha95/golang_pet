package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURLShorten(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte("http://vfdfbdfbd.com")))
	request.Header.Add("Content-type", "text/plain")

	w := httptest.NewRecorder()
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost: `http://localhost:8080`,
		FilePath: ``,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}
	h := GetSaveURLShortenHandler(stCfg)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, echo.MIMEApplicationJSON, w.Header().Get(echo.HeaderContentType))
	_, err := url.ParseRequestURI(w.Body.String())
	require.NoError(t, err)

	resBody, _ := io.ReadAll(res.Body)

	_, err = url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)
}

func TestInvalidURLShorten(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte("this is not URL")))
	request.Header.Add("Content-type", "text/plain")

	w := httptest.NewRecorder()
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost: `http://localhost:8080`,
		FilePath: ``,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}
	h := GetSaveURLShortenHandler(stCfg)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode)
}