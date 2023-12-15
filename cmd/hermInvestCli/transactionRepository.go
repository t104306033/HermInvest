package main

import (
	"database/sql"
	"fmt"
)

type transactionRepository struct {
	db *sql.DB
}

// CreateTransaction: insert transaction and return inserted id
func (repo *transactionRepository) createTransaction(t *Transaction) (int, error) {
	query := `INSERT INTO tblTransaction (stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := repo.db.Exec(query, t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
	if err != nil {
		fmt.Println("Error insert database: ", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting inserted id: ", err)
		return 0, err
	}
	return int(id), nil
}
