package repository

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/Vaha95/golang_pet/internal/config"
	dto "github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
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

func (s DBStrategy) SaveBatch(data []dto.BatchItem, generateHash func() string) error {    
	valueStrings := make([]string, 0, len(data))
    valueArgs := make([]interface{}, 0, len(data) * 3)
    for k, batch := range data {
		id := generateHash()
		i := (k+1)*3
        valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i-2, i-1, i))
        valueArgs = append(valueArgs, id, batch.URL, batch.ExtId)
    }

    stmt := fmt.Sprintf("INSERT INTO url_short (short,url,ext_id) VALUES %s ON CONFLICT (url) DO NOTHING", 
                        strings.Join(valueStrings, ","))
    _, err := s.db.Exec(stmt, valueArgs...)

    for k, batch := range data {
		short, err := s.getShortByURL(batch.URL)
		if err != nil {
			return fmt.Errorf("failed to find short: %w", err)
		}

		data[k].Short = short	
	}

    return err
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

func (s DBStrategy) getShortByURL(u string) (string, error) {
	id, err := s.GetShortByURL(u)
	if err != nil {
		return "", fmt.Errorf("failed to find short by URL: %w", err)
	}

	shortURL, err := url.JoinPath(s.cfg.URLHost, id)
	if err != nil {
		return "", fmt.Errorf("failed to create URL from short: %w", err)
	}

	return shortURL, nil
}