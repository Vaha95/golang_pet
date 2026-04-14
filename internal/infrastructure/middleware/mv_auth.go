package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

const (
	SECRET_KEY = "SECRET_KEY"
	TOKEN_EXP = time.Hour * 3
	COOKIE_KEY = "auth_token"
	COOKIE_EXP = time.Hour * 24
)

type Claims struct {
    jwt.RegisteredClaims
    UserID int
}

func addAuthMiddleware(cfg config.StorageConfig, e *echo.Echo, l *zap.SugaredLogger) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if cfg.IsDBAllowed {
				token, err := readCookie(c)

				userId := getUserID(token)
				if userId < 0 || err != nil {
					userId, err = repository.CreateUser(cfg.DBService.GetDB())
					if err != nil {
						l.Errorf("buildJWTString err: %w", err)

						return err
					}

					token, err := buildJWTString()
					if err != nil {
						l.Errorf("buildJWTString err: %w", err)

						return err
					}
					writeCookie(c, token)
				}
				cfg.SetUserId(userId)				
			}

			next(c)

			return nil
		}
	})
}

func buildJWTString() (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims {
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
        },
        UserID: 1,
    })

    tokenString, err := token.SignedString([]byte(SECRET_KEY))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func getUserID(tokenString string) int {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims,
    func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(SECRET_KEY), nil
    })
    if err != nil {
        return -1
    }

    if !token.Valid {
        return -1
    }

    return claims.UserID
}

func writeCookie(c *echo.Context, token string) {
	cookie := new(http.Cookie)

	cookie.Name = COOKIE_KEY
	cookie.Value = token
	cookie.Expires = time.Now().Add(COOKIE_EXP)
	cookie.Path = "/"

	c.SetCookie(cookie)
}

func readCookie(c *echo.Context) (string, error) {
	cookie, err := c.Cookie(COOKIE_KEY)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}