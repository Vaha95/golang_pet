package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURL(t *testing.T) {
	res, auditCh := sendDefRequest(t, "http://vfdfbdfbd.com")
	assert.Equal(t, 201, res.StatusCode)

	resBody, _ := io.ReadAll(res.Body)

	_, err := url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)

	select {
		case auditItem := <-auditCh:
			assert.Equal(t, "http://vfdfbdfbd.com", auditItem.URL)
			assert.Equal(t, DTO.FOLLOW, auditItem.Action)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for audit item")
	}
}

func TestInvalidURL(t *testing.T) {
	res, _ := sendDefRequest(t, "this is not URL")

	assert.Equal(t, 400, res.StatusCode)
}

func sendDefRequest(t *testing.T, url string) (*http.Response, chan DTO.BaseAuditItem) {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte(url)))
	request.Header.Add("Content-type", "text/plain")

	w := httptest.NewRecorder()
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   ``,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}
	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1)
	h := GetSaveURLHandler(stCfg, l, auditCh)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()

	return res, auditCh
}
