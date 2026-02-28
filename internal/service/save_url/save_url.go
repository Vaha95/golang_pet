package saveurl

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
)

func SaveURL(cfg config.Config, inputURL string) (string, error) {
		parsedURL, err := url.ParseRequestURI(inputURL)
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrorParseRequestURI, parsedURL)
		}

		id, err := setToStorage(cfg, parsedURL.String())
		if err != nil {
			return "", fmt.Errorf("%w: %s", ErrorSaveToStorage, parsedURL)
		}

	return url.JoinPath(cfg.URLHost, id)

}

func setToStorage(cfg config.Config, parsedURL string) (string, error) {
	for i := 0; i < 10; i++ {
		id := GenerateHash()
		if err := repository.SetURL(cfg, id, parsedURL); err != nil {
			if errors.Is(err, repository.ErrorShortURLKeyAlreadyExists) {
				continue
			}
			return "", fmt.Errorf("failed to save the short URL: %w", err)
		}

		return id, nil
	}

	return "", errors.New("ID generate is impossible")
}

func GenerateHash() string {
	rand.New((rand.NewSource(time.Now().UnixNano())))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}