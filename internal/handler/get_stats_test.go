package handler

import (
	"database/sql"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/stats"
	_ "github.com/jackc/pgx/v5/stdlib" // регистрация драйвера pgx
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockStatsDB struct {
	db *sql.DB
}

func (m *mockStatsDB) Close() error   { return nil }
func (m *mockStatsDB) Ping() error    { return nil }
func (m *mockStatsDB) GetDB() *sql.DB { return m.db }

func makeStatsCfg(db *sql.DB, trustedSubnet string) config.StorageConfig {
	return config.StorageConfig{
		DBService:   &mockStatsDB{db: db},
		IsDBAllowed: db != nil,
		Config: config.Config{
			TrustedSubnet: trustedSubnet,
		},
	}
}

func TestStats_NoTrustedSubnet_ReturnsForbidden(t *testing.T) {
	l, _ := getLogger()
	cfg := makeStatsCfg(nil, "")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "10.0.0.5")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestStats_TrustedIP_ReturnsOK(t *testing.T) {
	l, _ := getLogger()
	cfg := makeStatsCfg(nil, "10.0.0.0/8")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "10.0.0.5")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestStats_UntrustedIP_ReturnsForbidden(t *testing.T) {
	l, _ := getLogger()
	cfg := makeStatsCfg(nil, "10.0.0.0/8")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "192.168.1.1")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestStats_NoXRealIP_ReturnsForbidden(t *testing.T) {
	l, _ := getLogger()
	cfg := makeStatsCfg(nil, "10.0.0.0/8")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestStats_InvalidCIDR_ReturnsForbidden(t *testing.T) {
	l, _ := getLogger()
	cfg := makeStatsCfg(nil, "not-a-cidr")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "10.0.0.5")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestStats_DBError_Returns500(t *testing.T) {
	l, _ := getLogger()
	db, err := sql.Open("pgx", "invalid_dsn")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	cfg := makeStatsCfg(db, "0.0.0.0/0")
	h := GetStatsHandler(cfg, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "1.2.3.4")
	w := httptest.NewRecorder()
	c := echo.New().NewContext(req, w)
	h(c)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestIsInSubnet(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
	}{
		{"ipv4 inside", "10.0.0.5", "10.0.0.0/8", true},
		{"ipv4 outside", "192.168.0.1", "10.0.0.0/8", false},
		{"exact boundary", "10.0.0.0", "10.0.0.0/8", true},
		{"invalid cidr", "10.0.0.5", "bad", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			got := isInSubnet(ip, tt.cidr)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCollector_NoDB_ReturnsEmpty(t *testing.T) {
	c := stats.New(nil, false)
	resp, err := c.Collect()
	assert.NoError(t, err)
	assert.Equal(t, stats.Response{}, resp)
}
