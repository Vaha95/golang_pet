package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	echo "github.com/labstack/echo/v4"
)

func TestSaveUrl(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte("http://vfdfbdfbd.com")))
	request.Header.Add("Content-type", "text/plain")

	storage := repository.NewStorage()

	w := httptest.NewRecorder()
	urlHost := `localhost:8080`
	h := GetSaveURLHandler(storage, &urlHost)

	c := echo.New().NewContext(request, w)
	h(c)

	res := w.Result()
	assert.Equal(t, 201, res.StatusCode)

	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	_, err := url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)
}

func TestGenerateID(t *testing.T) {
	id := generateHash()
	require.Equal(t, 8, len(id))
}