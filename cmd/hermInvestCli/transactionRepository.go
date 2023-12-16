package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type transactionRepository struct {
	db *sql.DB
}

func (repo *transactionRepository) prepareStmt(sqlStmt string, tx *sql.Tx) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error

	if tx == nil {
		stmt, err = repo.db.Prepare(sqlStmt)
	} else {
		stmt, err = tx.Prepare(sqlStmt)
	}

	return stmt, err
}

func (repo *transactionRepository) createTransactionWithTx(t *Transaction, tx *sql.Tx) (int, error) {
	const insertSql string = "" +
		"INSERT INTO tblTransaction" +
		"(stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	stmt, err := repo.prepareStmt(insertSql, tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
	if err != nil {
		fmt.Println("Error insert database: ", err)
		return 0, err
	}

	id, err := rst.LastInsertId()
	if err != nil {
		fmt.Println("Error getting inserted id: ", err)
		return 0, err
	}

	return int(id), nil
}

// CreateTransaction: insert transaction and return inserted id
func (repo *transactionRepository) createTransaction(t *Transaction) (int, error) {
	return repo.createTransactionWithTx(t, nil)
}

// testcase begin, commit, rollback
// CreateTransactions: insert transactions and return inserted ids
func (repo *transactionRepository) createTransactions(ts []*Transaction) ([]int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var insertedIDs []int
	for _, t := range ts {
		id, err := repo.createTransactionWithTx(t, tx)
		if err != nil {
			fmt.Println("Error create transaction with Tx: ", err)
			return nil, err
		}
		insertedIDs = append(insertedIDs, int(id))
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return insertedIDs, nil
}

// queryTransactionAll
func (repo *transactionRepository) queryTransactionAll() ([]*Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// queryTransactionByID
func (repo *transactionRepository) queryTransactionByID(id int) ([]*Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var transactions []*Transaction
	var t Transaction
	err := row.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, &t)

	return transactions, nil
}

// queryTransactionByDetails
func (repo *transactionRepository) queryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error) {
	var conditions []string
	var args []interface{}

	if stockNo != "" {
		conditions = append(conditions, "stockNo = ?")
		args = append(args, stockNo)
	}
	if tranType != 0 {
		conditions = append(conditions, "tranType = ?")
		args = append(args, tranType)
	}
	if date != "" {
		conditions = append(conditions, "date = ?")
		args = append(args, date)
	}

	query := fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.id, &t.stockNo, &t.tranType, &t.quantity, &t.unitPrice, &t.totalAmount, &t.taxes)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// updateTransaction
func (repo *transactionRepository) updateTransaction(t *Transaction) error {
	_, err := repo.db.Exec("UPDATE tblTransaction SET unitPrice = ?, totalAmount = ?, taxes = ? WHERE id = ?", t.unitPrice, t.totalAmount, t.taxes, t.id)
	return err
}

// deleteTransaction
func (repo *transactionRepository) deleteTransaction(id int) error {
	_, err := repo.db.Exec("DELETE FROM tblTransaction WHERE id = ?", id)
	return err
}
