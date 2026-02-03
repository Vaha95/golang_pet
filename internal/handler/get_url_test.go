package handler

import (
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrl(t *testing.T) {
	url := "http://vfdfbdfbd.com"
	id := generateID()
	data := make(map[string]string)
	data[id] = url

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
	request.Header.Add("Content-type", "text/plain")
		vars := map[string]string{
		"id": id,
	}

	request = mux.SetURLVars(request, vars)

	w := httptest.NewRecorder()

	h := GetGetURLHandler(&data)

	h(w, request)

	res := w.Result()
	assert.Equal(t, 307, res.StatusCode)

	defer res.Body.Close()

	loc := res.Header.Get("Location")
	require.Equal(t, url, loc)
}