package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	dto "github.com/Vaha95/golang_pet/internal/model/DTO"
)

// DBStrategy implements URLStrategy using PostgreSQL as the backend.
type DBStrategy struct {
	db  *sql.DB
	cfg config.Config
}

// ErrorUrlAlreadyExists is returned when saving a URL that was already shortened.
var ErrorUrlAlreadyExists = errors.New("short URL already exists")

// ErrorUrlNotFound is returned when the URL doesn't exist in the database.
var ErrorUrlNotFound = errors.New("URL not found")

func (s DBStrategy) Save(short string, url string, extId *string, userId *int) error {
	query := "INSERT INTO url_short (url, short, ext_id, created_by_user) VALUES ($1,$2,$3,$4) ON CONFLICT (url) WHERE (deleted_at IS NULL) DO NOTHING RETURNING short"
	var res string
	err := s.db.QueryRow(query, url, short, extId, *userId).Scan(&res)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find short by URL: %w", ErrorUrlAlreadyExists)
		}

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

	stmt := fmt.Sprintf("INSERT INTO url_short (short,url, ext_id, created_by_user) VALUES %s ON CONFLICT (url) WHERE (deleted_at IS NULL) DO NOTHING RETURNING short, url",
		strings.Join(valueStrings, ","))

	rows, err := s.db.Query(stmt, valueArgs...)
	if err != nil {
		return err
	}

	shortMap := make(map[string]string, len(data))
	for rows.Next() {
		var short, u string
		if err := rows.Scan(&short, &u); err != nil {
			rows.Close()
			return err
		}
		shortMap[u] = short
	}
	rows.Close()

	for k, batch := range data {
		short, ok := shortMap[batch.URL]
		if !ok {
			short, err = s.GetShortByURL(batch.URL)
		}
		if err != nil {
			return fmt.Errorf("failed to find short: %w", err)
		}
		short, joinErr := url.JoinPath(s.cfg.URLHost, short)
		if joinErr != nil {
			return fmt.Errorf("failed to create URL from short: %w", joinErr)
		}

		data[k].Short = short
	}

	return nil
}

func (s DBStrategy) Get(short string) (*DTO.ShortItem, error) {
	stmt := "SELECT url, deleted_at FROM url_short where short=$1 LIMIT 1"

	row := s.db.QueryRow(stmt, short)

	var data DTO.ShortItem
	err := row.Scan(&data.URL, &data.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find URL: %w", ErrorUrlNotFound)
		}

		return nil, fmt.Errorf("failed to find URL: %w", err)
	}

	return &data, nil
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

func (s DBStrategy) GetByUser(userId *int) ([]DTO.ShortItem, error) {
	stmt := "SELECT short, url FROM url_short where created_by_user=$1"

	rows, err := s.db.Query(stmt, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find URL: %w", ErrorUrlNotFound)
		}

		return make([]dto.ShortItem, 0), fmt.Errorf("failed to parse data from DBRow: %w", err)
	}
	defer rows.Close()

	rows.Columns()

	var data []DTO.ShortItem
	for rows.Next() {
		var s DTO.ShortItem
		if err := rows.Scan(&s.Short, &s.URL); err != nil {
			return make([]dto.ShortItem, 0), err
		}
		data = append(data, s)
	}

	return data, nil
}

func (s DBStrategy) DeleteBatch(data DTO.DeleteBatch) error {
	batch := data.Shorts

	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*2)

	for k, short := range batch {
		i := (k + 1) * 2
		valueStrings = append(valueStrings, fmt.Sprintf("short = $%d AND created_by_user=$%d AND deleted_at IS NULL", i-1, i))
		valueArgs = append(valueArgs, short, data.UserId)
	}

	stmt := fmt.Sprintf("UPDATE url_short SET deleted_at = now() WHERE (%s)",
		strings.Join(valueStrings, ") OR ("))
	_, err := s.db.Exec(stmt, valueArgs...)

	if err != nil {
		return err
	}

	return err
}
