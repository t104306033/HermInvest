package repository

import (
	"HermInvest/pkg/model"
	"database/sql"
	"fmt"
	"strings"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
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

func (repo *transactionRepository) createTransactionWithTx(t *model.Transaction, tx *sql.Tx) (int, error) {
	const insertSql string = "" +
		"INSERT INTO tblTransaction" +
		"(stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	stmt, err := repo.prepareStmt(insertSql, tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(t.StockNo, t.Date, t.Quantity, t.TranType, t.UnitPrice, t.TotalAmount, t.Taxes)
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
func (repo *transactionRepository) CreateTransaction(t *model.Transaction) (int, error) {
	return repo.createTransactionWithTx(t, nil)
}

// testcase begin, commit, rollback
// CreateTransactions: insert transactions and return inserted ids
func (repo *transactionRepository) CreateTransactions(ts []*model.Transaction) ([]int, error) {
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
func (repo *transactionRepository) QueryTransactionAll() ([]*model.Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
func (repo *transactionRepository) QueryTransactionByID(id int) ([]*model.Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var transactions []*model.Transaction
	var t model.Transaction
	err := row.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, &t)

	return transactions, nil
}

// queryTransactionByDetails
func (repo *transactionRepository) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
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

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
func (repo *transactionRepository) UpdateTransaction(t *model.Transaction) error {
	_, err := repo.db.Exec("UPDATE tblTransaction SET unitPrice = ?, totalAmount = ?, taxes = ? WHERE id = ?", t.UnitPrice, t.TotalAmount, t.Taxes, t.ID)
	return err
}

// deleteTransaction
func (repo *transactionRepository) DeleteTransaction(id int) error {
	_, err := repo.db.Exec("DELETE FROM tblTransaction WHERE id = ?", id)
	return err
}
