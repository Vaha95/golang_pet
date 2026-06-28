package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockDBService struct {
	pingErr error
}

func (m *mockDBService) Ping() error { return m.pingErr }

func TestPingDB_Success(t *testing.T) {
	l, _ := getLogger()
	h := GetPingDBHandler(&mockDBService{pingErr: nil}, l)

	w := httptest.NewRecorder()
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "http://localhost/ping", nil), w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPingDB_Fail(t *testing.T) {
	l, _ := getLogger()
	h := GetPingDBHandler(&mockDBService{pingErr: assert.AnError}, l)

	w := httptest.NewRecorder()
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "http://localhost/ping", nil), w)
	h(c)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
