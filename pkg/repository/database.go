package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection() (*gorm.DB, error) {
	// TODO: extract DB connection to configuration
	// sqlite3 connection with foreign keys enabled
	// TODO: need to check db file, if not exist, exit. otherwise ... db will be created
	var dbPath = "./internal/app/database/dev-database.db"

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Reference: https://gorm.io/docs/logger.html
		// GORM defined log levels: Silent (default), Error, Warn, Info.
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
