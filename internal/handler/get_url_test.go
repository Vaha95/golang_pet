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

func BenchmarkTest(b *testing.B) {
}
