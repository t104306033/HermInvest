package repository

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBConnection() *gorm.DB {
	// TODO: extract DB connection to configuration
	// TODO: TestCase check db file, if not exist, exit.
	var dbPath = "./internal/app/database/dev-database.db"

	err := isFileExist(dbPath)
	if err != nil {
		fmt.Println("[Error] Datebase path is invalid:", err)
		fmt.Println("\n* Please make sure the database file exists.")
		os.Exit(1)
	}

	// Reference: https://gorm.io/docs/logger.html
	// GORM defined log levels: Silent (default), Error, Warn, Info.
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		// Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		fmt.Println("[Error] Failed to open database:", err)
		os.Exit(1)
	}

	// Reference: https://stackoverflow.com/questions/3888529/how-to-tell-if-sqlite-database-file-is-valid-or-not
	err = isSQLiteFile(db)
	if err != nil {
		fmt.Println("[Error] Unable to use database:", err)
		fmt.Printf("\n* The database path is: '%s'.\n", dbPath)
		os.Exit(1)
	}

	err = checkDBSchema(db)
	if err != nil {
		fmt.Println("[Error] Unable to ensure database schema integrity:", err)
		fmt.Println("\n* Please check schema or consider rebuilding the database.")
		os.Exit(1)
	}

	return db
}

func isFileExist(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file '%s' does not exist: %w", path, errors.Unwrap(err))
		}
		return fmt.Errorf("unable to obtain the info of path '%s': %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("path '%s' is a directory not a file", path)
	}

	return nil
}

func isSQLiteFile(db *gorm.DB) error {
	var schemaVersion int
	err := db.Raw("PRAGMA schema_version;").Scan(&schemaVersion).Error
	if err != nil {
		return fmt.Errorf("failed to execute 'PRAGMA schema_version': %w", err)
	}

	if schemaVersion == 0 {
		return fmt.Errorf("SQLite database is empty or uninitialized (schema version: %d)", schemaVersion)
	}

	return nil
}

func checkDBSchema(db *gorm.DB) error {
	var result string
	err := db.Raw("PRAGMA integrity_check;").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("failed to execute 'PRAGMA integrity_check': %w", err)
		// Or use 'PRAGMA schema.quick_check;'
	}

	if result != "ok" {
		return fmt.Errorf("schema check failed (check result: %s)", result)
	}

	return nil
}
