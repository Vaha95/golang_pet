package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestDeleteURL_Success(t *testing.T) {
	ch := make(chan DTO.DeleteBatch, 1)
	stCfg := config.StorageConfig{}
	stCfg.SetUserId(42)

	l, _ := getLogger()
	h := GetDeleteURLHandler(stCfg, ch, l)

	payload, _ := json.Marshal([]string{"abc123", "def456"})
	request := httptest.NewRequest(http.MethodPost, "http://localhost/delete", bytes.NewReader(payload))
	request.Header.Add("Content-type", "application/json")

	w := httptest.NewRecorder()
	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)

	msg := <-ch
	assert.Equal(t, 42, msg.UserId)
	assert.Equal(t, []string{"abc123", "def456"}, msg.Shorts)
}

func TestDeleteURL_EmptyArray(t *testing.T) {
	ch := make(chan DTO.DeleteBatch, 1)
	stCfg := config.StorageConfig{}
	stCfg.SetUserId(1)

	l, _ := getLogger()
	h := GetDeleteURLHandler(stCfg, ch, l)

	payload := []byte(`[]`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost/delete", bytes.NewReader(payload))
	request.Header.Add("Content-type", "application/json")

	w := httptest.NewRecorder()
	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteURL_InvalidJSON(t *testing.T) {
	ch := make(chan DTO.DeleteBatch, 1)
	stCfg := config.StorageConfig{}
	stCfg.SetUserId(1)

	l, _ := getLogger()
	h := GetDeleteURLHandler(stCfg, ch, l)

	request := httptest.NewRequest(http.MethodPost, "http://localhost/delete", bytes.NewReader([]byte(`not json`)))
	request.Header.Add("Content-type", "application/json")

	w := httptest.NewRecorder()
	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
