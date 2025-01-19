// internal/database/sqlite.go
package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func NewSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "policies.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}