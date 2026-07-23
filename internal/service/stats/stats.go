package stats

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vaha95/golang_pet/internal/config"
)

var (
	ErrorCountURLs  = errors.New("failed to count URLs")
	ErrorCountUsers = errors.New("failed to count users")
)

// Response is the JSON body of GET /api/internal/stats.
type Response struct {
	URLsInDB int `json:"urls"`
	Users    int `json:"users"`
}

// Collector gathers stats from a database.
type Collector struct {
	db    *sql.DB
	alive bool
}

// New returns a Collector. If db is nil or not alive, queries are skipped.
func New(db *sql.DB, alive bool) *Collector {
	return &Collector{db: db, alive: alive}
}

// Collect returns current internal stats.
func (c *Collector) Collect() (Response, error) {
	if !c.alive || c.db == nil {
		return Response{}, nil
	}

	var totalURLs, users int

	if err := c.db.QueryRow("SELECT count(1) FROM url_short WHERE deleted_at IS NULL").Scan(&totalURLs); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrorCountURLs, err)
	}

	if err := c.db.QueryRow("SELECT count(1) FROM users").Scan(&users); err != nil {
		return Response{}, fmt.Errorf("%w: %w", ErrorCountUsers, err)
	}

	return Response{
		URLsInDB: totalURLs,
		Users:    users,
	}, nil
}

// CollectFromConfig collects stats from a StorageConfig. It is a convenience wrapper for the handler.
func CollectFromConfig(cfg config.StorageConfig) (Response, error) {
	var db *sql.DB
	if cfg.DBService != nil {
		db = cfg.DBService.GetDB()
	}
	return New(db, cfg.IsDBAllowed).Collect()
}
