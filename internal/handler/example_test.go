package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v5"

	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

// TestExampleSaveURL demonstrates saving a raw URL from the request body.
func ExampleSaveURL() {
	cfg := config.StorageConfig{
		Config: config.Config{
			ListenHost: "localhost:8080",
			URLHost:    "http://localhost:8080",
		},
	}
	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1)
	h := GetSaveURLHandler(cfg, l, auditCh)

	req := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:8080/",
		bytes.NewReader([]byte("https://example.com")),
	)
	req.Header.Set("Content-Type", "text/plain")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	h(c)

	fmt.Printf("status: %d\n", rec.Code)
	fmt.Printf("body: %s\n", rec.Body.String())
}

// TestExampleGetSaveURLShortenHandler demonstrates saving a URL as JSON.
func ExampleGetSaveURLShortenHandler() {
	cfg := config.StorageConfig{
		Config: config.Config{
			ListenHost: "localhost:8080",
			URLHost:    "http://localhost:8080",
		},
	}
	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1)
	h := GetSaveURLShortenHandler(cfg, l, auditCh)

	payload, _ := json.Marshal(map[string]string{
		"url": "https://example.com",
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"http://localhost:8080/api/shorten",
		bytes.NewReader(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	h(c)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)

	fmt.Printf("status: %d\n", rec.Code)
	fmt.Printf("result: %s\n", resp["result"])
}

// TestExampleGetURLHandler demonstrates resolving a short URL and receiving a redirect.
func ExampleGetURLHandler() {
	cfg := config.StorageConfig{
		Config: config.Config{
			ListenHost: "localhost:8080",
			URLHost:    "http://localhost:8080",
		},
	}
	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1)
	h := GetURLHandler(cfg, l, auditCh)

	originalURL := "https://example.com"
	id := saveurl.GenerateHash()

	// Pre-seed the URL in storage.
	repository.SetURL(cfg.Config, id, originalURL)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:8080/%s", id),
		nil,
	)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetPathValues(echo.PathValues{
		{Name: "id", Value: id},
	})
	h(c)

	fmt.Printf("status: %d\n", rec.Code)
	fmt.Printf("location: %s\n", rec.Header().Get("Location"))
}

// mockOKDB is a test double that always returns nil from Ping.
type mockDB struct{}

func (m *mockDB) Ping() error { return nil }

// TestExampleGetPingDBHandler demonstrates checking database connectivity.
func ExampleGetPingDBHandler() {
	l, _ := getLogger()

	h := GetPingDBHandler(&mockDB{}, l)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	h(c)

	fmt.Printf("status: %d\n", rec.Code)
}
