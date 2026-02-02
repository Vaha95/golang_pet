package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler(t *testing.T) {
	type want struct {
		code 	 	int
		contentType string
	}

	tests := []struct{
		name string
		want want
	}{
		{
			name: "pos test 12",
			want: want{
				code: 201,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal("test.com")
			if err != nil {
				t.Fatal(err)
			}

			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", bytes.NewReader(bodyBytes))
			request.Header.Add("Content-type", "text/plain")

			w := httptest.NewRecorder()
			saveUrl(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			_, err = url.Parse(string(resBody))
			if err != nil {
				t.Error(err.Error())
			}
			require.NoError(t, err)
		})
	}
}