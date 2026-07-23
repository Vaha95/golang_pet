package grpcserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authHeader = "authorization"

func AuthInterceptor(cfg *config.StorageConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userId, err := extractUserID(ctx, cfg)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		cfg.SetUserId(userId)

		return handler(ctx, req)
	}
}

func extractUserID(ctx context.Context, cfg *config.StorageConfig) (int, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return createNewUser(cfg)
	}

	tokenStr := extractToken(md)
	if tokenStr == "" {
		return createNewUser(cfg)
	}

	userId := getUserIdFromToken(tokenStr)
	if userId >= 0 {
		return userId, nil
	}

	return createNewUser(cfg)
}

func extractToken(md metadata.MD) string {
	vals := md.Get(authHeader)
	if len(vals) == 0 {
		return ""
	}

	val := vals[0]
	if strings.HasPrefix(val, "Bearer ") {
		return strings.TrimPrefix(val, "Bearer ")
	}

	return val
}

func getUserIdFromToken(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}

func createNewUser(cfg *config.StorageConfig) (int, error) {
	if cfg.IsDBAllowed && cfg.DBService != nil {
		userId, err := repository.CreateUser(cfg.DBService.GetDB())
		if err != nil {
			return -1, fmt.Errorf("failed to create user: %w", err)
		}
		return userId, nil
	}
	return 0, nil
}

// JWT claims shared with HTTP auth middleware.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// SecretKey used to sign and verify JWT tokens.
const SecretKey = "SECRET_KEY"
