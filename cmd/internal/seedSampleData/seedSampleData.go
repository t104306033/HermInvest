package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./internal/app/database/dev-database.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Insert tblStockMapping data
	_, err = db.Exec(`
		INSERT INTO tblStockMapping (stockNo, stockName) VALUES 
		('0050', 'tw50'),
		('0051', 'tw51'),
		('0052', 'tw52')
	`)
	if err != nil {
		fmt.Println(err)
	}

	// Insert tblTransaction data
	_, err = db.Exec(`
		INSERT INTO tblTransaction (id, stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES 
		(1, '0050', '2023-12-04', 1000, 1, 20, 20, 3),
		(2, '0050', '2023-12-05', 1000, 1, 20, 20, 3),
		(3, '0050', '2023-12-06', 1000, 1, 20, 20, 3),
		(4, '0051', '2023-12-06', 1000, 1, 20, 20, 3),
		(5, '0051', '2023-12-06', 1000, 1, 20, 20, 3),
		(6, '0051', '2023-12-06', 1000, 1, 20, 20, 3),
		(7, '0052', '2023-12-06', 1000, 1, 20, 20, 3),
		(8, '0052', '2023-12-06', 1000, 1, 20, 20, 3),
		(9, '0052', '2023-12-06', 1000, 1, 20, 20, 3),
		(10, '0052', '2023-12-06', 1000, 1, 20, 20, 3)
	`)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Seed data inserted successfully")
}
