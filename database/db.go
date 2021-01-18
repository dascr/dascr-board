package database

import (
	"database/sql"

	// Import of the sqlite3 driver
	"github.com/dascr/dascr-board/config"
	_ "github.com/mattn/go-sqlite3"
)

// SetupDB will instantiate a database and return it
func SetupDB(cfg *config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Filename)
	if err != nil {
		return nil, err
	}

	// Table Generation statement
	create := `
	CREATE TABLE IF NOT EXISTS player (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		nickname TEXT,
		image TEXT
	);
	`

	// Create tables
	_, err = db.Exec(create)
	if err != nil {
		return nil, err
	}

	// Return DB
	return db, nil
}
