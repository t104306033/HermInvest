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

// ---
// Transaction

// createTransactionWithTx: insert transaction and return inserted id
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

// ---
// Transaction History

// createTransactionHistoryWithTx: insert transaction and return inserted id
func (repo *transactionRepository) createTransactionHistoryWithTx(t *model.Transaction, tx *sql.Tx) (int, error) {
	const insertSql string = "" +
		"INSERT INTO tblTransactionHistory" +
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

// CreateTransactionHistory: insert transaction and return inserted id
func (repo *transactionRepository) CreateTransactionHistory(t *model.Transaction) (int, error) {
	return repo.createTransactionHistoryWithTx(t, nil)
}

// testcase begin, commit, rollback
// CreateTransactionHistorys: insert transactions and return inserted ids
func (repo *transactionRepository) CreateTransactionHistorys(ts []*model.Transaction) ([]int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var insertedIDs []int
	for _, t := range ts {
		id, err := repo.createTransactionHistoryWithTx(t, tx)
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

// ---

// FindFirstPurchase
func (repo *transactionRepository) FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error) {
	query := "" +
		"SELECT id, stockNo, date, tranType, quantity, unitPrice, totalAmount, taxes " +
		"FROM tblTransaction WHERE stockNo = ? " +
		"ORDER BY date ASC LIMIT 1"
	row := repo.db.QueryRow(query, stockNo)

	var t model.Transaction
	err := row.Scan(&t.ID, &t.StockNo, &t.Date, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
	if err != nil {
		return &model.Transaction{}, err
	}

	// fmt.Println(t.ID, t.StockNo, t.Date, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes)

	return &t, nil
}

// QueryInventoryTransactions
func (repo *transactionRepository) QueryInventoryTransactions(stockNo string, quantity int) ([]*model.Transaction, error) {
	query := "" +
		"WITH cte AS (" +
		"	SELECT *, SUM(quantity) OVER (ORDER BY date, id) AS running_total" +
		"	FROM tblTransaction" +
		"	WHERE stockNo = ?" +
		") " +
		"SELECT id, stockNo, date, tranType, quantity, unitPrice, totalAmount, taxes " +
		"FROM cte WHERE running_total <= ?"

	rows, err := repo.db.Query(query, stockNo, quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.StockNo, &t.Date, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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

// queryTransactionAll
func (repo *transactionRepository) QueryTransactionAll() ([]*model.Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, date, unitPrice, totalAmount, taxes FROM tblTransaction`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.Date, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
func (repo *transactionRepository) QueryTransactionByID(id int) (*model.Transaction, error) {
	query := `SELECT id, stockNo, tranType, quantity, date, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var t model.Transaction
	err := row.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.Date, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
	if err != nil {
		return nil, err
	}

	return &t, nil
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

	query := fmt.Sprintf("SELECT id, stockNo, tranType, quantity, date, unitPrice, totalAmount, taxes FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.StockNo, &t.TranType, &t.Quantity, &t.Date, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
func (repo *transactionRepository) UpdateTransaction(id int, t *model.Transaction) error {
	query := "" +
		"UPDATE tblTransaction " +
		"SET stockNo = ?, date = ?, quantity = ?, tranType = ?, unitPrice = ?, totalAmount = ?, taxes = ? " +
		"WHERE id = ?"
	_, err := repo.db.Exec(query, t.StockNo, t.Date, t.Quantity, t.TranType, t.UnitPrice, t.TotalAmount, t.Taxes, t.ID)
	return err
}

// deleteTransaction
func (repo *transactionRepository) DeleteTransaction(id int) error {
	_, err := repo.db.Exec("DELETE FROM tblTransaction WHERE id = ?", id)
	return err
}

// deleteTransactions
func (repo *transactionRepository) DeleteTransactions(ids []int) error {
	var args []interface{}
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
		args = append(args, ids[i])
	}

	query := fmt.Sprintf("DELETE FROM tblTransaction WHERE id IN (%s)", strings.Join(placeholders, ", "))

	_, err := repo.db.Exec(query, args...)
	return err

}

// MoveInventoryToTransactionHistorys
func (repo *transactionRepository) MoveInventoryToTransactionHistorys(ts []*model.Transaction) error {
	// TODO: operate SQL with TX
	var ids []int
	for _, t := range ts {
		ids = append(ids, t.ID)
	}

	// Delete transactions from transaction
	err := repo.DeleteTransactions(ids)
	if err != nil {
		return fmt.Errorf("error deleting transactions: %v", err)
	}

	// Create transactions to from transactionHistory
	_, err = repo.CreateTransactionHistorys(ts)
	if err != nil {
		return fmt.Errorf("error creating transaction history: %v", err)
	}

	return nil
}

// AddTransaction add the transaction from the input to the inventory.
// It will add or update transactions in the inventory and add history.
// Return the modified transaction record in the inventory
func (repo *transactionRepository) AddTransaction(newTransaction *model.Transaction) ([]*model.Transaction, error) {
	// TODO: This func should be moved to service tier.

	// Find the first purchase from the inventory
	remainingQuantity := newTransaction.Quantity
	earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
	if err == sql.ErrNoRows {
		// If no first purchase found, set TranType and handle accordingly
		earliestTransaction.TranType = newTransaction.TranType
	} else if err != nil {
		return nil, fmt.Errorf("error finding first purchase: %v", err)
	}

	if earliestTransaction.TranType == newTransaction.TranType {
		// add stock in the inventory directly
		id, err := repo.CreateTransaction(newTransaction)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %v", err)
		}
		transactions, err := repo.QueryTransactionByID(id)
		if err != nil {
			return nil, fmt.Errorf("error querying database: %v", err)
		}

		return transactions, nil
	}

	inventoryTransactions, err := repo.QueryInventoryTransactions(newTransaction.StockNo, newTransaction.Quantity)
	if err != nil {
		return nil, fmt.Errorf("error querying old purchase: %v", err)
	}

	for _, t := range inventoryTransactions {
		remainingQuantity -= t.Quantity
	}

	err = repo.MoveInventoryToTransactionHistorys(inventoryTransactions)
	if err != nil {
		return nil, fmt.Errorf("error moving inventory to transaction history: %v", err)
	}

	if remainingQuantity > 0 {
		// Find the first purchase again after old purchases are processed
		earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("add stock in the inventory directly", err)

				// record newTransaction to history
				newTransaction.SetQuantity(newTransaction.Quantity - remainingQuantity)
				_, err = repo.CreateTransactionHistory(newTransaction)
				if err != nil {
					// handle create transaction history failed
				}

				newTransaction.SetQuantity(remainingQuantity)
				id, err := repo.CreateTransaction(newTransaction)
				if err != nil {
					return nil, fmt.Errorf("error creating transaction: %v", err)
				}

				transactions, err := repo.QueryTransactionByID(id)
				if err != nil {
					return nil, fmt.Errorf("error querying database: %v", err)
				}

				return transactions, nil

			}

			return nil, fmt.Errorf("error querying database: %v", err)
		}

		if earliestTransaction.Quantity < remainingQuantity {
			fmt.Println("Not Impelment yet: earliestTransaction.Quantity < remainingQuantity: ", err)
			return nil, fmt.Errorf("not Impelment yet: earliestTransaction.Quantity < remainingQuantity: %v", err)
		}

		// Update the quantity of the first purchase
		earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
		err = repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
		if err != nil {
			// handle update transaction failed
		}

		// record purchase history
		earliestTransaction.SetQuantity(remainingQuantity)
		_, err = repo.CreateTransactionHistory(earliestTransaction)
		if err != nil {
			// handle create transaction history failed
		}
	}

	// record newTransaction to history
	_, err = repo.CreateTransactionHistory(newTransaction)
	if err != nil {
		// handle create transaction history failed
	}

	return nil, nil
}
