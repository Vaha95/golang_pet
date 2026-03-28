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
	Host string
	Port int
	User string
	Password string
	DBName string
	SSLMode string // например, "disable" или "require"

	MaxOpenConns int // максимальное количество открытых соединений
	MaxIdleConns int // максимальное количество простаивающих соединений
	ConnMaxLifetime time.Duration // максимальное время жизни соединения
}

// Connect устанавливает соединение с базой данных
func connect(cfg Config) (*sql.DB, error) {
	db, err := ConnectToDB(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func getDSN(cfg Config) string {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)		
	}

	return dsn
}

func ConnectToDB(cfg Config) (*sql.DB, error) {
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

func InitDb() *Service {
	mainConfig := config.GetMainConfig()
	cfg := Config{
		DSN: mainConfig.DbDSN,
		Host: "localhost",
		Port: 5432,
		User: "myuser",
		Password: "mypass",
		DBName: "mydatabase",
		SSLMode: "disable",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := connect(cfg)
	if ; err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	
    err = initMigrations(db)
	if err != nil {
		log.Printf("Failed to migrate: %v", err)
	}

	return NewService(db)
}