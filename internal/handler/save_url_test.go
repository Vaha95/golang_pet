package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveUrl(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader([]byte("http://vfdfbdfbd.com")))
	request.Header.Add("Content-type", "text/plain")

	data := make(map[string]string)

	w := httptest.NewRecorder()
	h := GetSaveURLHandler(&data)

	h(w, request)

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
	id := generateID()
	require.Equal(t, 8, len(id))
}