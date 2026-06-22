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
	strategy "github.com/Vaha95/golang_pet/internal/repository/strategy/save_url"
)

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
