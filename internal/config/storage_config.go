package config

import (
	"database/sql"
)

// StorageConfig bundles the storage backend, flag for DB availability, app config, and user data.
type StorageConfig struct {
	DBService   DBService
	IsDBAllowed bool
	Config      Config
	UserData    *UserData
}

// DBService provides database operations required by the storage layer.
type DBService interface {
	Close() error
	Ping() error
	GetDB() *sql.DB
}

// GetConfig creates a StorageConfig from the given database service, flag, and app config.
func GetConfig(dbService DBService, isDBAllowed bool, mainConfig Config) StorageConfig {
	return StorageConfig{
		DBService:   dbService,
		IsDBAllowed: isDBAllowed,
		Config:      mainConfig,
		UserData:    &UserData{},
	}
}

// GetUserId returns the authenticated user ID, or nil if no user is set.
func (cfg *StorageConfig) GetUserId() *int {
	if cfg.UserData == nil {
		return nil
	}

	return &(cfg.UserData.userId)
}

// SetUserId sets the authenticated user ID in the config.
func (cfg *StorageConfig) SetUserId(userId int) {
	if cfg.UserData == nil {
		cfg.UserData = &UserData{userId}

		return
	}

	cfg.UserData.userId = userId
}
