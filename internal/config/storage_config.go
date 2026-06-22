package config

import (
	"database/sql"
)

type StorageConfig struct {
	DBService   DBService
	IsDBAllowed bool
	Config      Config
	UserData    *UserData
}

type DBService interface {
	Close() error
	Ping() error
	GetDB() *sql.DB
}

func GetConfig(dbService DBService, isDBAllowed bool, mainConfig Config) StorageConfig {
	return StorageConfig{
		DBService:   dbService,
		IsDBAllowed: isDBAllowed,
		Config:      mainConfig,
		UserData:    &UserData{},
	}
}

func (cfg *StorageConfig) GetUserId() *int {
	if cfg.UserData == nil {
		return nil
	}

	return &(cfg.UserData.userId)
}

func (cfg *StorageConfig) SetUserId(userId int) {
	if cfg.UserData == nil {
		cfg.UserData = &UserData{userId}
		
		return
	}

	cfg.UserData.userId = userId
}
