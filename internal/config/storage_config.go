package config

import (
	"database/sql"
)

type StorageConfig struct {
	DBService DBService
	IsDBAllowed bool
	Config Config
	UserData *UserData
}

type DBService interface {
	Close() error
	Ping() error
	GetDB() *sql.DB
}

func GetConfig(dbService DBService, isDBAllowed bool, mainConfig Config) StorageConfig {
	return StorageConfig{
		DBService: dbService,
		IsDBAllowed: isDBAllowed,
		Config: mainConfig,
		UserData: &UserData{},
	}
}

func (cfg *StorageConfig) GetUserId() int {
	return cfg.UserData.userId
}

func (cfg *StorageConfig) SetUserId(userId int) {
	cfg.UserData.userId = userId
}