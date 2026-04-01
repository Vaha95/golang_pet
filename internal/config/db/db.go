package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // регистрация драйвера pgx
)

// Config содержит параметры подключения к PostgreSQL
type Config struct {
	DSN string

	MaxOpenConns int // максимальное количество открытых соединений
	MaxIdleConns int // максимальное количество простаивающих соединений
	ConnMaxLifetime time.Duration // максимальное время жизни соединения
}

// Connect устанавливает соединение с базой данных
func connect(cfg Config) (*sql.DB, error) {
	db, err := сonnectToDB(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func getDSN(cfg Config) string {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = "host=localhost port=5432 user=myuser password=mypass dbname=mydatabase sslmode=disable"
	}

	return dsn
}

func сonnectToDB(cfg Config) (*sql.DB, error) {
	dsn := getDSN(cfg)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

func InitDb(mainConfig config.Config) (*Service, error) {
	cfg := Config{
		DSN: mainConfig.DbDSN,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := connect(cfg)
	if ; err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}
	
    err = initMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to migrate: %v", err)
	}

	return NewService(db), nil
}