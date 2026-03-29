package repository

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/Vaha95/golang_pet/internal/config"
	DTO "github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
)

type DBStrategy struct {
	db *sql.DB
	cfg config.Config
}

func (s DBStrategy) Save(short string, url string, extId *string) error {
	sql := "INSERT INTO url_short (url,short,ext_id) VALUES ($1,$2,$3) ON CONFLICT (url) DO NOTHING"
	_, err := s.db.Exec(sql, url, short, extId)

	if err != nil {
		return err
	}

	return nil
}

func (s DBStrategy) SaveBatch(data []DTO.BatchItem, GenerateHash func() string) error {
	t, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction fo save to DB: %w", err)
	}
	defer t.Commit()

	for i := 0; i < len(data); i++ {
		item := &data[i]

		id := GenerateHash()
		extId := item.ExtId

		err := s.Save(id, item.URL, &extId)
		if err != nil {
			t.Rollback()

			return fmt.Errorf("failed to save the short URL to DB: %w", err)
		}

		id, err = s.GetShortByURL(item.URL)
		if err != nil {
			t.Rollback()

			return fmt.Errorf("failed to find short by URL: %w", err)
		}

		shortURL, err := url.JoinPath(s.cfg.URLHost, id)
		if err != nil {
			t.Rollback()

			return fmt.Errorf("failed to create URL from short: %w", err)
		}
		item.Short = shortURL
	}

	return nil
}

func (s DBStrategy) Get(short string) (string, error) {
	sql := "SELECT url FROM url_short where short=$1 LIMIT 1"

	row := s.db.QueryRow(sql, short)

	var url string
	err := row.Scan(&url)
	if err != nil || url == "" {
		return "", fmt.Errorf("failed to parse URL from DBRow: %w", err)
	}

	return url, nil
}

func (s DBStrategy) GetShortByURL(url string) (string, error) {
	sql := "SELECT short FROM url_short where url=$1 LIMIT 1"

	row := s.db.QueryRow(sql, url)

	var short string
	err := row.Scan(&short)
	if err != nil || short == "" {
		return "", fmt.Errorf("failed to parse short from DBRow: %w", err)
	}

	return short, nil
}