package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Vaha95/golang_pet/internal/config"
)

var ErrorShortURLKeyAlreadyExists = errors.New("short URL key already exists")
var ErrorShortURLKeyNotFound = errors.New("short URL key already exists")

func GetURLByKey(cfg config.Config, key string) (string, error) {
	data, err := ReadFileStore(cfg)
	if err != nil {
		return "", err
	}
	if data == nil {
		data = map[string]string{}
	}
	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrorShortURLKeyNotFound, key)
	}

	return val, nil
}

func SetURL(cfg config.Config, key string, val string) (err error) {
	data, err := ReadFileStore(cfg)
	if err != nil {
		return err
	}

	_, ok := data[key]
	if ok {
		return fmt.Errorf("%w: %s", ErrorShortURLKeyAlreadyExists, key)
	}
	if data == nil {
		data = map[string]string{}
	}

	data[key] = val
	WriteFileStore(cfg, data)

	return nil
}

func SetURLToDB(db *sql.DB, ctx context.Context, short string, url string, extId string) (err error) {
	sql := "INSERT INTO url_short (url,short,ext_id) VALUES ($1,$2,$3)"
	_, err = db.ExecContext(ctx, sql, url, short, extId)

	log.Println(sql)

	if err != nil {
		return err
	}

	return nil
}