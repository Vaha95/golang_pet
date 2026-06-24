package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"go.uber.org/zap"

	mv "github.com/Vaha95/golang_pet/internal/infrastructure/middleware"
)

type APIReqiest struct {
	URI string `json:"url"`
}

type APIResponse struct {
	Result string `json:"result"`
}

func TestSaveURLShorten(t *testing.T) {
	res, w, auditCh := sendDefShortenRequest(t, "http://bgfnfgmnhg.com")

	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, echo.MIMEApplicationJSON, w.Header().Get(echo.HeaderContentType))

	var resData APIResponse
	json.Unmarshal(w.Body.Bytes(), &resData)
	_, err := url.ParseRequestURI(resData.Result)
	require.NoError(t, err)

	select {
		case auditItem := <-auditCh:
			assert.Equal(t, "http://bgfnfgmnhg.com", auditItem.URL)
			assert.Equal(t, DTO.FOLLOW, auditItem.Action)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for audit item")
	}
}

func TestInvalidURLShorten(t *testing.T) {
	res, _, _ := sendDefShortenRequest(t, "this is not URL")

	assert.Equal(t, 400, res.StatusCode)
}

func sendDefShortenRequest(t *testing.T, url string) (*http.Response, *httptest.ResponseRecorder, chan DTO.BaseAuditItem) {
	data, err := json.Marshal(APIReqiest{url})
	require.NoError(t, err)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewReader(data))
	request.Header.Add("Content-type", "application/json")

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
	h := GetSaveURLShortenHandler(stCfg, l, auditCh)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()

	return res, w, auditCh
}

func getLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("%w", mv.ErrorZapLoggerInitialize)
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}
