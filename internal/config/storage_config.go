package config

import (
	"database/sql"
)

type StorageConfig struct {
	DBService DBService
	IsDBAllowed bool
	Config Config
}

type DBService interface {
	Close() error
	Ping() error
	GetDB() *sql.DB
}

func GetConfig(dbService DBService, isDBAllowed bool) StorageConfig {
	return StorageConfig{
		DBService: dbService,
		IsDBAllowed: isDBAllowed,
		Config: GetMainConfig(),
	}
}