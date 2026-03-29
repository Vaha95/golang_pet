package repository

import (
	"database/sql"
	"fmt"

	DTO "github.com/Vaha95/golang_pet/internal/model/DTO/save_url"
)

type DBStrategy struct {
	db *sql.DB
}

func (s DBStrategy) Save(short string, url string, extId *string) error {
	sql := "INSERT INTO url_short (url,short,ext_id) VALUES ($1,$2,$3)"
	_, err := s.db.Exec(sql, url, short, extId)

	if err != nil {
		return err
	}

	return nil
}

func (s DBStrategy) SaveBatch(data *[]DTO.BatchItem, GenerateHash func() string) error {
	t, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction fo save to DB: %w", err)
	}
	defer t.Commit()

	for i := 0; i < len(*data); i++ {
		id := GenerateHash()
		extId := (*data)[i].ExtId
		err := s.Save(id, (*data)[i].URL, &extId)
		if err != nil {
			t.Rollback()

			return fmt.Errorf("failed to save the short URL to DB: %w", err)
		}

		(*data)[i].Short = id
	}

	return nil
}