package db

import (
	"database/sql"
	"errors"
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	
	path, err := migrationPath()
	if err != nil {
		return err
	}

    m, err := migrate.NewWithDatabaseInstance(
        "file://" + path,
        "postgres",
		driver,
	)
	
	if (err != nil) {
		return err
	}

    err = m.Up() 
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func migrationPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b) 
	basepath, _ = url.JoinPath(basepath, "migrations")
	
	return basepath, nil
}