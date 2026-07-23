package grpcserver

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Vaha95/golang_pet/api/shortenerpb"
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	saveurl "github.com/Vaha95/golang_pet/internal/service/save_url"
)

func newTestServer(t *testing.T) (*Server, chan DTO.BaseAuditItem) {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "store.json")
	auditCh := make(chan DTO.BaseAuditItem, 1)
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := config.StorageConfig{Config: cfg}
	return NewServer(stCfg, auditCh), auditCh
}

// --- Shorten ---

func TestShorten_Success(t *testing.T) {
	srv, auditCh := newTestServer(t)
	resp, err := srv.Shorten(context.Background(), &shortenerpb.ShortenRequest{
		Url: "http://example.com",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.ShortUrl)

	parsedURL, err := url.ParseRequestURI(resp.ShortUrl)
	require.NoError(t, err)
	assert.Equal(t, "http", parsedURL.Scheme)

	assertAudit(t, auditCh, "http://example.com")
}

func TestShorten_InvalidURL(t *testing.T) {
	srv, _ := newTestServer(t)
	_, err := srv.Shorten(context.Background(), &shortenerpb.ShortenRequest{
		Url: "not a valid url",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestShorten_EmptyURL(t *testing.T) {
	srv, _ := newTestServer(t)
	_, err := srv.Shorten(context.Background(), &shortenerpb.ShortenRequest{
		Url: "",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestShorten_ValidHTTPSURL(t *testing.T) {
	srv, auditCh := newTestServer(t)
	resp, err := srv.Shorten(context.Background(), &shortenerpb.ShortenRequest{
		Url: "https://secure.example.com/path?query=1",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.ShortUrl)
	assertAudit(t, auditCh, "https://secure.example.com/path?query=1")
}

// --- Resolve ---

func TestResolve_Success(t *testing.T) {
	srv, auditCh := newTestServer(t)
	id := seedURL(srv.cfg.Config, "http://resolve.com")

	resp, err := srv.Resolve(context.Background(), &shortenerpb.ResolveRequest{Id: id})
	require.NoError(t, err)
	assert.Equal(t, "http://resolve.com", resp.Url)
	assert.False(t, resp.Deleted)

	assertAudit(t, auditCh, "http://resolve.com")
}

func TestResolve_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	_, err := srv.Resolve(context.Background(), &shortenerpb.ResolveRequest{
		Id: "nonexistent_id",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestResolve_EmptyID(t *testing.T) {
	srv, _ := newTestServer(t)
	_, err := srv.Resolve(context.Background(), &shortenerpb.ResolveRequest{
		Id: "",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

// --- GetUserURLs ---

func TestGetUserURLs_Success(t *testing.T) {
	srv, _ := newTestServer(t)

	repository.SetURL(srv.cfg.Config, "short1", "http://url1.com")
	repository.SetURL(srv.cfg.Config, "short2", "http://url2.com")

	resp, err := srv.GetUserURLs(context.Background(), &shortenerpb.GetUserURLsRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Urls, 2)

	got := make(map[string]string, len(resp.Urls))
	for _, u := range resp.Urls {
		got[u.Short] = u.Url
	}
	assert.Equal(t, "http://url1.com", got["short1"])
	assert.Equal(t, "http://url2.com", got["short2"])
}

func TestGetUserURLs_Empty(t *testing.T) {
	srv, _ := newTestServer(t)
	resp, err := srv.GetUserURLs(context.Background(), &shortenerpb.GetUserURLsRequest{})

	require.NoError(t, err)
	require.Empty(t, resp.Urls)
}

func TestGetUserURLs_MixedWithShorten(t *testing.T) {
	srv, _ := newTestServer(t)

	_, err := srv.Shorten(context.Background(), &shortenerpb.ShortenRequest{
		Url: "http://mixed.com",
	})
	require.NoError(t, err)

	resp, err := srv.GetUserURLs(context.Background(), &shortenerpb.GetUserURLsRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Urls, 1)
	assert.Equal(t, "http://mixed.com", resp.Urls[0].Url)
}

// --- benchmarks ---

func BenchmarkShorten(b *testing.B) {
	auditCh := make(chan DTO.BaseAuditItem, b.N)
	tmpFile := filepath.Join(os.TempDir(), "bench_shorten.json")
	defer os.Remove(tmpFile)

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := config.StorageConfig{Config: cfg}
	srv := NewServer(stCfg, auditCh)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = srv.Shorten(ctx, &shortenerpb.ShortenRequest{
			Url: "http://benchmark.com",
		})
	}
}

func BenchmarkResolve(b *testing.B) {
	tmpFile := filepath.Join(os.TempDir(), "bench_resolve.json")
	defer os.Remove(tmpFile)

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := config.StorageConfig{Config: cfg}
	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, "http://benchmark.com")

	auditCh := make(chan DTO.BaseAuditItem, b.N)
	srv := NewServer(stCfg, auditCh)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = srv.Resolve(ctx, &shortenerpb.ResolveRequest{Id: id})
	}
}

func BenchmarkResolveNotFound(b *testing.B) {
	tmpFile := filepath.Join(os.TempDir(), "bench_resolve_nf.json")
	defer os.Remove(tmpFile)

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := config.StorageConfig{Config: cfg}

	auditCh := make(chan DTO.BaseAuditItem, 1024)
	srv := NewServer(stCfg, auditCh)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = srv.Resolve(ctx, &shortenerpb.ResolveRequest{Id: "no_such_id"})
	}
}

func BenchmarkGetUserURLs(b *testing.B) {
	tmpFile := filepath.Join(os.TempDir(), "bench_user_urls.json")
	defer os.Remove(tmpFile)

	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := config.StorageConfig{Config: cfg}

	for i := 0; i < 50; i++ {
		repository.SetURL(cfg, "bench_short"+string(rune(i)), "http://bench.com/"+string(rune(i)))
	}

	auditCh := make(chan DTO.BaseAuditItem, 1024)
	srv := NewServer(stCfg, auditCh)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = srv.GetUserURLs(ctx, &shortenerpb.GetUserURLsRequest{})
	}
}

// --- helpers ---

func seedURL(cfg config.Config, targetURL string) string {
	id := saveurl.GenerateHash()
	repository.SetURL(cfg, id, targetURL)
	return id
}

func assertAudit(t *testing.T, ch chan DTO.BaseAuditItem, expectedURL string) {
	t.Helper()
	select {
	case item := <-ch:
		assert.Equal(t, expectedURL, item.URL)
		assert.Equal(t, DTO.FOLLOW, item.Action)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for audit item")
	}
}
