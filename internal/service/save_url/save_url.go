package saveurl

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/repository"
	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

var (
	hashChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	hashRand  = rand.New(rand.NewSource(time.Now().UnixNano()))
	hashMu    sync.Mutex
)

// SaveURL validates and persists a single URL, returning the full short URL path.
func SaveURL(cfg config.StorageConfig, inputURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrorParseRequestURI, parsedURL)
	}

	id, err := setToStorage(cfg, parsedURL.String())
	if err != nil {
		if errors.Is(err, ErrorUrlAlreadyExists) {
			id, joinErr := url.JoinPath(cfg.Config.URLHost, id)
			if joinErr != nil {
				return "", joinErr
			}

			return id, err
		} else {
			return "", fmt.Errorf("%w: %s, %w", ErrorSaveToStorage, parsedURL, err)
		}
	}

	return url.JoinPath(cfg.Config.URLHost, id)

}

func setToStorage(cfg config.StorageConfig, parsedURL string) (string, error) {
	for i := 0; i < 10; i++ {
		id := GenerateHash()
		strg := strategy.GetStrategy(cfg)
		if err := strg.Save(id, parsedURL, nil, cfg.GetUserId()); err != nil {
			if errors.Is(err, repository.ErrorShortURLKeyAlreadyExists) {
				continue
			}

			if errors.Is(err, strategy.ErrorUrlAlreadyExists) {
				id, err = strg.GetShortByURL(parsedURL)
				if err != nil {
					return id, fmt.Errorf("failed to find short by URL: %w", err)
				}

				return id, fmt.Errorf("failed to save new URL: %w", ErrorUrlAlreadyExists)
			}

			return "", fmt.Errorf("failed to save the short URL: %w", err)
		}

		return id, nil
	}

	return "", errors.New("ID generate is impossible")
}

// GenerateHash produces a random 8-character alphanumeric short URL key.
func GenerateHash() string {
	hashMu.Lock()
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = hashChars[hashRand.Intn(len(hashChars))]
	}
	hashMu.Unlock()

	return string(buf)
}
