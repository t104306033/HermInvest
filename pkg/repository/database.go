package repository

import (
	"database/sql"
	"fmt"
)

func GetDBConnection() (*sql.DB, error) {
	// TODO: extract DB connection to configuration
	// sqlite3 connection with foreign keys enabled
	// TODO: need to check db file, if not exist, exit. otherwise ... db will be created
	var dbPath = "./internal/app/database/dev-database.db"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return db, nil
}
