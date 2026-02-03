package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveUrl(t *testing.T) {
	bodyBytes, err := json.Marshal("test.com")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(bodyBytes))
	request.Header.Add("Content-type", "text/plain")

	w := httptest.NewRecorder()
	saveUrl(w, request)

	res := w.Result()
	assert.Equal(t, 201, res.StatusCode)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	_, err = url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)
}

func TestGetUrl(t *testing.T) {
	bodyBytes, err := json.Marshal("test.com")
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(bodyBytes))
	request.Header.Add("Content-type", "text/plain")

	w := httptest.NewRecorder()
	saveUrl(w, request)

	request = httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	request.Header.Add("Content-type", "text/plain")
	mux.SetURLVars(request, map[string]string{"id": "test"})

	w = httptest.NewRecorder()
	getUrl(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	fmt.Println(string(resBody))
	assert.Equal(t, 201, res.StatusCode)
	_, err = url.Parse(string(resBody))
	if err != nil {
		t.Error(err.Error())
	}
	require.NoError(t, err)
}