package main

import (
	"database/sql"
	"fmt"
)

func GetDBConnection() (*sql.DB, error) {
	// TODO: extract DB connection to configuration
	// sqlite3 connection with foreign keys enabled
	var dbPath = "./internal/app/database/dev-database.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return db, nil
}
