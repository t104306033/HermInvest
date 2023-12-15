package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create SQLite Database
	db, err := sql.Open("sqlite3", "./internal/app/database/dev-database.db")
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}
	defer db.Close()
	fmt.Println("Database is creating ...")

	// Create tblStockMapping table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblStockMapping (
			stockNo TEXT NOT NULL UNIQUE,
			stockName TEXT NOT NULL,
			PRIMARY KEY(stockNo)
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblStockMapping table:", err)
		return
	}
	fmt.Println("Table tblStockMapping created successfully")

	// Create tblTransaction table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tblTransaction (
			id INTEGER UNIQUE,
			stockNo TEXT NOT NULL,
			date TEXT,
			quantity INTEGER NOT NULL,
			tranType INTEGER NOT NULL,
			unitPrice REAL NOT NULL,
			totalAmount INTEGER,
			taxes INTEGER,
			PRIMARY KEY("id" AUTOINCREMENT)
		)
	`)
	if err != nil {
		fmt.Println("Error creating tblTransaction table:", err)
		return
	}
	fmt.Println("Table tblTransaction created successfully")

	fmt.Println("Database created successfully")
}
