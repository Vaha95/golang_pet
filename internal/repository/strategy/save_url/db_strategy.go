package repository

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	dto "github.com/Vaha95/golang_pet/internal/model/DTO"
)

type DBStrategy struct {
	db  *sql.DB
	cfg config.Config
}

func (s DBStrategy) Save(short string, url string, extId *string, userId *int) error {
	sql := "INSERT INTO url_short (url, short, ext_id, created_by_user) VALUES ($1,$2,$3,$4) ON CONFLICT (url) DO NOTHING"
	_, err := s.db.Exec(sql, url, short, extId, *userId)

	if err != nil {
		return err
	}

	return nil
}

func (s DBStrategy) SaveBatch(data []dto.BatchItem, userId *int, generateHash func() string) error {
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*4)
	for k, batch := range data {
		id := generateHash()
		i := (k + 1) * 4
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i-3, i-2, i-1, i))
		valueArgs = append(valueArgs, id, batch.URL, batch.ExtId, userId)
	}

	stmt := fmt.Sprintf("INSERT INTO url_short (short,url, ext_id, created_by_user) VALUES %s ON CONFLICT (url) DO NOTHING",
		strings.Join(valueStrings, ","))
	_, err := s.db.Exec(stmt, valueArgs...)

	for k, batch := range data {
		short, err := s.GetShortByURL(batch.URL)
		if err != nil {
			return fmt.Errorf("failed to find short: %w", err)
		}
		short, err = url.JoinPath(s.cfg.URLHost, short)
		if err != nil {
			return fmt.Errorf("failed to create URL from short: %w", err)
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

func (s DBStrategy) GetShortByURL(u string) (string, error) {
	sql := "SELECT short FROM url_short where url=$1 LIMIT 1"

	row := s.db.QueryRow(sql, u)

	var short string
	err := row.Scan(&short)
	if err != nil || short == "" {
		return "", fmt.Errorf("failed to parse short from DBRow: %w", err)
	}

	return short, nil
}

func (s DBStrategy) GetByUser(userId int) ([]DTO.ShortItem, error) {
	sql := "SELECT short, url FROM url_short where created_by_user=$1 LIMIT 1"

	row := s.db.QueryRow(sql, userId)

	var data []DTO.ShortItem
	err := row.Scan(&data)
	if err != nil || len(data) <= 0 {
		return make([]dto.ShortItem, 0), fmt.Errorf("failed to parse data from DBRow: %w", err)
	}

	return data, nil
}