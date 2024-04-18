package repository

import (
	"HermInvest/pkg/model"
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

type TransactionRepositoryGorm struct {
	db *gorm.DB
}

func NewTransactionRepositoryGorm(db *gorm.DB) *TransactionRepositoryGorm {
	return &TransactionRepositoryGorm{db: db}
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
		"(date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := repo.prepareStmt(insertSql, tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(t.Date, t.Time, t.StockNo, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes)
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
		"(date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	stmt, err := repo.prepareStmt(insertSql, tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rst, err := stmt.Exec(t.Date, t.Time, t.StockNo, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes)
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
		"SELECT id, date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes " +
		"FROM tblTransaction WHERE stockNo = ? " +
		"ORDER BY date ASC LIMIT 1"
	row := repo.db.QueryRow(query, stockNo)

	var t model.Transaction
	err := row.Scan(&t.ID, &t.Date, &t.Time, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
		"SELECT id, date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes " +
		"FROM cte WHERE running_total <= ?"

	rows, err := repo.db.Query(query, stockNo, quantity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Time, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
	query := `SELECT id, date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Time, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
	query := `SELECT id, date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var t model.Transaction
	err := row.Scan(&t.ID, &t.Date, &t.Time, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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

	query := fmt.Sprintf("SELECT id, date, time, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.ID, &t.Date, &t.Time, &t.StockNo, &t.TranType, &t.Quantity, &t.UnitPrice, &t.TotalAmount, &t.Taxes)
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
		"SET date = ?, time = ?, stockNo = ?, tranType = ?, quantity = ?, unitPrice = ?, totalAmount = ?, taxes = ? " +
		"WHERE id = ?"
	_, err := repo.db.Exec(query, t.Date, t.Time, t.StockNo, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes, t.ID)
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

func (repo *TransactionRepositoryGorm) QueryCapitalReductionAll() ([]model.CapitalReduction, error) {
	var capitalReductions []model.CapitalReduction
	// 使用 Gorm 框架的 Find 方法來執行查詢
	if err := repo.db.Debug().Table("tblCapitalReduction").Find(&capitalReductions).Error; err != nil {
		return nil, err
	}
	fmt.Println(capitalReductions)
	fmt.Println(capitalReductions[0])
	return capitalReductions, nil
}

// queryTransactionAll
func (repo *TransactionRepositoryGorm) QueryTransactionRecordByStockNo(stockNo string) ([]model.TransactionRecord, error) {
	var transactionRecords []model.TransactionRecord
	err := repo.db.Debug().Table("tblTransactionRecord").Where("stockNo = ?", stockNo).Find(&transactionRecords).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(transactionRecords)
	return transactionRecords, nil
}

// Service Tier

// addTransactionTailRecursion add new transaction records with tail recursion,
// When adding, inventory and transaction history, especially write-offs and
// tails, need to be considered.
func (repo *transactionRepository) addTransactionTailRecursion(newTransaction *model.Transaction, remainingQuantity int) (*model.Transaction, error) {
	// Principles:
	// 1. Ensure that each transaction has a corresponding transaction record.
	// 2. Update inventory quantities based on transactions, including adding,
	//    reducing, or deleting inventory.
	// 3. Depending on the transaction situation, only transaction history can
	//    be added and cannot be modified or deleted.
	// 4. For insufficient write-off quantities, recursive processing is used
	//    to ensure that the write-off is completed.

	// Cases:
	// 1. Newly added: If there is no transaction in the inventory (A) or
	//    the new transaction is the same as the oldest transaction in the
	//    inventory (B), add it directly to the inventory.
	// 2. Write-off:
	// 	* Sufficient inventory: If the inventory quantity is sufficient,
	//    update the inventory quantity (C) or delete the inventory (D), and
	//    add the corresponding transaction history.
	// 	* Insufficient inventory: If the inventory quantity can't be Write-off.
	//    Recurse until success (E). The termination condition is A B C D.
	//  * Over inventory: Write-off over than inventory (F).

	// TODO: This func should be moved to service tier.

	earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("error finding first purchase: %v", err)
		}
		// Case A
		earliestTransaction.TranType = newTransaction.TranType
	}

	if earliestTransaction.TranType == newTransaction.TranType {
		if newTransaction.Quantity != remainingQuantity {
			// Case F
			newTransaction.SetQuantity(newTransaction.Quantity - remainingQuantity)
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			newTransaction.SetQuantity(remainingQuantity)
		}

		// Case B
		id, err := repo.CreateTransaction(newTransaction)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %v", err)
		}
		transaction, err := repo.QueryTransactionByID(id)
		if err != nil {
			return nil, fmt.Errorf("error querying database: %v", err)
		}

		return transaction, nil
	} else {
		if earliestTransaction.Quantity > remainingQuantity {
			// Case C

			// Create a copy for adding stock history
			stockHistoryAdd := &model.Transaction{}
			*stockHistoryAdd = *earliestTransaction
			// var stockHistoryAdd *model.Transaction // why can't use it, study it
			// *stockHistoryAdd = *earliestTransaction

			// add transaction history
			stockHistoryAdd.SetQuantity(remainingQuantity)
			_, err = repo.CreateTransactionHistory(stockHistoryAdd)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// Update stock inventory
			earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
			err := repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
			if err != nil {
				// handle update transaction failed
			}

			return earliestTransaction, nil
		} else if earliestTransaction.Quantity == remainingQuantity {
			// Case D

			// add transaction history
			_, err = repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			// delete stock inventory
			err = repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			// Or use move

			return nil, nil
		} else { // earliestTransaction.Quantity < remainingQuantity
			// Case E

			// add transaction history
			_, err = repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// delete stock inventory
			err = repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			remainingQuantity = remainingQuantity - earliestTransaction.Quantity

			return repo.addTransactionTailRecursion(newTransaction, remainingQuantity)
		}
	}
}

// AddTransaction add the transaction from the input to the inventory.
// It will add or update transactions in the inventory and add history.
// Return the modified transaction record in the inventory
func (repo *transactionRepository) AddTransaction(newTransaction *model.Transaction) (*model.Transaction, error) {
	remainingQuantity := newTransaction.Quantity
	return repo.addTransactionTailRecursion(newTransaction, remainingQuantity)
}
