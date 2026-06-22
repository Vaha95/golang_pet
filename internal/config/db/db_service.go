package db

import (
	"database/sql"
	"fmt"
)

// Service представляет сервис для работы с PostgreSQL
type Service struct {
	db *sql.DB
}

// NewService создаёт новый экземпляр сервиса с заданной конфигурацией
func NewService(db *sql.DB) *Service {
	return &Service{db}
}

// Close закрывает соединение с базой данных
func (s *Service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// Ping проверяет доступность базы данных
func (s *Service) Ping() error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return s.db.Ping()
}

// GetDB возвращает объект *sql.DB для выполнения запросов
func (s *Service) GetDB() *sql.DB {
	return s.db
}
