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

func TestSaveURL(t *testing.T) {
	res := sendDefRequest("http://vfdfbdfbd.com")
	assert.Equal(t, 201, res.StatusCode)

	resBody, _ := io.ReadAll(res.Body)

	_, err := url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)
}

func TestInvalidURL(t *testing.T) {
	res := sendDefRequest("this is not URL")

	assert.Equal(t, 400, res.StatusCode)
}

func TestURLAlreadyExists(t *testing.T) {
	url := "http://bdngvvnbv.com"

	res := sendDefRequest(url)
	assert.Equal(t, 201, res.StatusCode)

	res = sendDefRequest(url)
	resBody, _ := io.ReadAll(res.Body)
	assert.Equal(t, 409, res.StatusCode, resBody)
}

func sendDefRequest(url string) *http.Response {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte(url)))
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
	h := GetSaveURLHandler(stCfg)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	return res
}