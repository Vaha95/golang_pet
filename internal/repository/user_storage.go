package repository

import (
	"database/sql"
	"fmt"
)

// CreateUser inserts a new row into the users table and returns the generated ID.
func CreateUser(db *sql.DB) (int, error) {
	sql := "INSERT INTO users DEFAULT VALUES RETURNING id;"

	row := db.QueryRow(sql)

	var userId int
	err := row.Scan(&userId)
	if err != nil || userId <= 0 {
		return 0, fmt.Errorf("failed to parse URL from DBRow: %w", err)
	}

	return userId, nil
}
