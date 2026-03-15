package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Vaha95/golang_pet/internal/config"
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

// Service представляет сервис для работы с PostgreSQL
type Service struct {
	db *sql.DB
}

// NewService создаёт новый экземпляр сервиса с заданной конфигурацией
func NewService(cfg Config) *Service {
	return &Service{}
}

// Connect устанавливает соединение с базой данных
func (s *Service) Connect(cfg Config) error {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)		
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Проверяем доступность базы данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	log.Println("Connected to PostgreSQL successfully")
	return nil
}

// Close закрывает соединение с базой данных
func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// Ping проверяет доступность базы данных
func (s *Service) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return s.db.PingContext(ctx)
}

// GetDB возвращает объект *sql.DB для выполнения запросов
func (s *Service) GetDB() *sql.DB {
	return s.db
}

func InitDb(mainConfig config.Config) {
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

	pgService := NewService(cfg)
	if err := pgService.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgService.Close()
}