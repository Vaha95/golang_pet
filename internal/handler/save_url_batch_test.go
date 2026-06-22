package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestSaveURLBatch_Success(t *testing.T) {
	payload, _ := json.Marshal([]map[string]string{
		{"correlation_id": "abc", "original_url": "http://example.com"},
		{"correlation_id": "def", "original_url": "http://test.com"},
	})

	res := sendBatchRequest(payload)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var resp []map[string]string
	json.NewDecoder(res.Body).Decode(&resp)

	assert.Len(t, resp, 2)
	assert.Equal(t, "abc", resp[0]["correlation_id"])
	assert.NotEmpty(t, resp[0]["short_url"])
	assert.Equal(t, "def", resp[1]["correlation_id"])
	assert.NotEmpty(t, resp[1]["short_url"])
}

func TestSaveURLBatch_InvalidJSON(t *testing.T) {
	res := sendBatchRequest([]byte(`not json`))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestSaveURLBatch_EmptyArray(t *testing.T) {
	res := sendBatchRequest([]byte(`[]`))
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func sendBatchRequest(body []byte) *http.Response {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(body))
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
	h := GetSaveURLBatchHandler(stCfg, l)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	return res
}
