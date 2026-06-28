package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func TestGetUrl(t *testing.T) {
	url := "http://njknonnjkn.com"

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}

	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, url)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
	w := httptest.NewRecorder()
	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1)
	h := GetURLHandler(stCfg, l, auditCh)

	c := echo.New().NewContext(request, w)
	c.SetPathValues(echo.PathValues{
		{Name: "id", Value: id},
	})
	h(c)

	res := w.Result()
	assert.Equal(t, 307, res.StatusCode)

	defer res.Body.Close()

	loc := res.Header.Get("Location")
	require.Equal(t, url, loc)

	select {
	case auditItem := <-auditCh:
		assert.Equal(t, url, auditItem.URL)
		assert.Equal(t, DTO.FOLLOW, auditItem.Action)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for audit item")
	}
}

func BenchmarkGetURLHandler(b *testing.B) {
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}

	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, "http://njknonnjkn.com")

	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1024)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", id), nil)
		w := httptest.NewRecorder()
		h := GetURLHandler(stCfg, l, auditCh)

		c := echo.New().NewContext(request, w)
		c.SetPathValues(echo.PathValues{
			{Name: "id", Value: id},
		})
		h(c)

		w.Result().Body.Close()

		// drain audit channel to avoid blocking
		select {
		case <-auditCh:
		default:
		}
	}
}

func BenchmarkGetURLHandlerNotFound(b *testing.B) {
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}

	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, 1024)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/not_exist", nil)
		w := httptest.NewRecorder()
		h := GetURLHandler(stCfg, l, auditCh)

		c := echo.New().NewContext(request, w)
		c.SetPathValues(echo.PathValues{
			{Name: "id", Value: "not_exist"},
		})
		h(c)

		w.Result().Body.Close()
	}
}

func BenchmarkGetURLHandlerConcurrent(b *testing.B) {
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
	}
	stCfg := config.StorageConfig{
		Config: cfg,
	}

	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, "http://example.com")

	l, _ := getLogger()
	auditCh := make(chan DTO.BaseAuditItem, b.N)

	b.SetParallelism(8)
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

		for pb.Next() {
			w := httptest.NewRecorder()
			h := GetURLHandler(stCfg, l, auditCh)

			c := echo.New().NewContext(req, w)
			c.SetPathValues(echo.PathValues{
				{Name: "id", Value: id},
			})
			h(c)

			res := w.Result()
			res.Body.Close()

			select {
			case <-auditCh:
			default:
			}
		}
	})
}
