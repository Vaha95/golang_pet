package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/Vaha95/golang_pet/api/shortenerpb"
	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAuthInterceptor_ValidJWT(t *testing.T) {
	stCfg, userId, interceptor := setupInterceptor(t)
	token, _ := buildJWTForUser(userId)
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	_, err := interceptor(
		ctx,
		&shortenerpb.ShortenRequest{Url: "http://example.com"},
		&grpc.UnaryServerInfo{FullMethod: "/shortener.ShortenerService/Shorten"},
		func(_ context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, userId, *stCfg.GetUserId())
			return &shortenerpb.ShortenResponse{}, nil
		},
	)
	require.NoError(t, err)
}

func TestAuthInterceptor_NoAuth(t *testing.T) {
	_, _, interceptor := setupInterceptor(t)

	ctx := context.Background()

	_, err := interceptor(
		ctx,
		&shortenerpb.ShortenRequest{Url: "http://example.com"},
		&grpc.UnaryServerInfo{FullMethod: "/shortener.ShortenerService/Shorten"},
		func(_ context.Context, req interface{}) (interface{}, error) {
			return &shortenerpb.ShortenResponse{}, nil
		},
	)
	require.NoError(t, err)
}

func TestAuthInterceptor_ExpiredJWT(t *testing.T) {
	_, _, interceptor := setupInterceptor(t)

	token := buildExpiredJWT()
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer "+token),
	)

	_, err := interceptor(
		ctx,
		&shortenerpb.ShortenRequest{Url: "http://example.com"},
		&grpc.UnaryServerInfo{FullMethod: "/shortener.ShortenerService/Shorten"},
		func(_ context.Context, req interface{}) (interface{}, error) {
			return &shortenerpb.ShortenResponse{}, nil
		},
	)
	require.NoError(t, err)
}

func TestAuthInterceptor_InvalidToken(t *testing.T) {
	_, _, interceptor := setupInterceptor(t)

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer some.garbage.token"),
	)

	_, err := interceptor(
		ctx,
		&shortenerpb.ShortenRequest{Url: "http://example.com"},
		&grpc.UnaryServerInfo{FullMethod: "/shortener.ShortenerService/Shorten"},
		func(_ context.Context, req interface{}) (interface{}, error) {
			return &shortenerpb.ShortenResponse{}, nil
		},
	)
	require.NoError(t, err)
}

func TestAuthInterceptor_RawTokenNoBearer(t *testing.T) {
	stCfg, userId, interceptor := setupInterceptor(t)
	token, _ := buildJWTForUser(userId)
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", token),
	)

	_, err := interceptor(
		ctx,
		&shortenerpb.ShortenRequest{Url: "http://example.com"},
		&grpc.UnaryServerInfo{FullMethod: "/shortener.ShortenerService/Shorten"},
		func(_ context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, userId, *stCfg.GetUserId())
			return &shortenerpb.ShortenResponse{}, nil
		},
	)
	require.NoError(t, err)
}

// --- helpers ---

func setupInterceptor(t *testing.T) (*config.StorageConfig, int, grpc.UnaryServerInterceptor) {
	t.Helper()
	tmpFile := t.TempDir() + "/store.json"
	auditCh := make(chan DTO.BaseAuditItem, 1)
	cfg := config.Config{
		ListenHost: `localhost:8080`,
		URLHost:    `http://localhost:8080`,
		FilePath:   tmpFile,
	}
	stCfg := &config.StorageConfig{Config: cfg}
	_ = auditCh

	userId := 42
	interceptor := AuthInterceptor(stCfg)
	return stCfg, userId, interceptor
}

func buildJWTForUser(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserID: userId,
	})
	return token.SignedString([]byte(SecretKey))
}

func buildExpiredJWT() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
		UserID: 42,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return ""
	}
	return tokenString
}
