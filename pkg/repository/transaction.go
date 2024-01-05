package repository

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepositoryGorm struct {
	db *gorm.DB
}

func NewTransactionRepositoryGorm(db *gorm.DB) *TransactionRepositoryGorm {
	return &TransactionRepositoryGorm{db: db}
}

// ---
// Transaction

// CreateTransaction: insert transaction and return inserted id
func (repo *TransactionRepositoryGorm) CreateTransaction(t *model.Transaction) (int, error) {
	result := repo.db.Debug().Table("tblTransaction").Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// CreateTransactions: insert transactions and return inserted ids
func (repo *TransactionRepositoryGorm) CreateTransactions(ts []*model.Transaction) ([]int, error) {
	result := repo.db.Debug().Table("tblTransaction").Create(&ts)
	if result.Error != nil {
		return nil, result.Error
	}

	var insertedIDs []int
	for _, t := range ts {
		insertedIDs = append(insertedIDs, t.ID)
	}

	return insertedIDs, nil
}

// ---
// Transaction History

// CreateTransactionHistory: insert transaction and return inserted id
func (repo *TransactionRepositoryGorm) CreateTransactionHistory(t *model.Transaction) (int, error) {
	result := repo.db.Debug().Table("tblTransactionHistory").Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// testcase begin, commit, rollback
// CreateTransactionHistorys: insert transactions and return inserted ids
func (repo *TransactionRepositoryGorm) CreateTransactionHistorys(ts []*model.Transaction) ([]int, error) {
	result := repo.db.Debug().Table("tblTransactionHistory").Create(&ts)
	if result.Error != nil {
		return nil, result.Error
	}

	var insertedIDs []int
	for _, t := range ts {
		insertedIDs = append(insertedIDs, t.ID)
	}

	return insertedIDs, nil
}

// ---

func (repo *TransactionRepositoryGorm) FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Where("stockNo = ?", stockNo).
		Order("date ASC, time ASC").First(&transaction).Error
	if err != nil {
		return &model.Transaction{}, err
	}

	return &transaction, nil
}

// QueryTransactionAll
func (repo *TransactionRepositoryGorm) QueryTransactionAll() ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionByID
func (repo *TransactionRepositoryGorm) QueryTransactionByID(id int) (*model.Transaction, error) {
	var transaction *model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Take(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// QueryTransactionByDetails
func (repo *TransactionRepositoryGorm) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	if stockNo != "" {
		repo.db = repo.db.Where("stockNo = ?", stockNo)
	}
	if tranType != 0 {
		repo.db = repo.db.Where("tranType = ?", tranType)
	}
	if date != "" {
		repo.db = repo.db.Where("date = ?", date)
	}

	err := repo.db.Debug().Table("tblTransaction").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// updateTransaction
func (repo *TransactionRepositoryGorm) UpdateTransaction(id int, t *model.Transaction) error {
	err := repo.db.Debug().Table("tblTransaction").Updates(t).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepositoryGorm) DeleteTransaction(id int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{ID: id})
	return result.Error
}

func (repo *TransactionRepositoryGorm) DeleteTransactions(ids []int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{}, "id IN ?", ids)
	return result.Error
}

// QueryCapitalReductionAll
func (repo *TransactionRepositoryGorm) QueryCapitalReductionAll() ([]*model.CapitalReduction, error) {
	var capitalReductions []*model.CapitalReduction
	// 使用 Gorm 框架的 Find 方法來執行查詢
	if err := repo.db.Debug().Table("tblCapitalReduction").Find(&capitalReductions).Error; err != nil {
		return nil, err
	}

	// for _, cr := range capitalReductions {
	// 	fmt.Println(cr)
	// }

	return capitalReductions, nil
}

// QueryTransactionRecordByStockNo
func (repo *TransactionRepositoryGorm) QueryTransactionRecordByStockNo(stockNo string, date string) ([]*model.TransactionRecord, error) {
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Debug().Table("tblTransactionRecord").Where("stockNo = ? and date < ?", stockNo, date).Find(&transactionRecords).Error
	if err != nil {
		return nil, err
	}

	// for _, cr := range transactionRecords {
	// 	fmt.Println(cr)
	// }

	return transactionRecords, nil
}

// deleteAlltblTransaction
func (repo *TransactionRepositoryGorm) DeleteAlltblTransaction() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransaction").Error; err != nil {
		return err
	}

	return nil
}

// deleteAlltblTransactionHistory
func (repo *TransactionRepositoryGorm) DeleteAlltblTransactionHistory() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransactionHistory").Error; err != nil {
		return err
	}

	return nil
}

// deleteAllTransactionRecordSys
func (repo *TransactionRepositoryGorm) DeleteAllTransactionRecordSys() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransactionRecordSys").Error; err != nil {
		return err
	}

	return nil
}

// QueryUnionNote
func (repo *TransactionRepositoryGorm) QueryUnionNote() {
	// Most ORMs seem not support UNION keyword, due to its complexity.
	// Faced with this situation, community suggest using "Raw" method to do this.
	// > Reference:
	// > * https://github.com/go-gorm/gorm/issues/3781
	// > * https://stackoverflow.com/questions/67190972/how-to-use-mysql-union-all-on-gorm
	// > * https://gorm.io/docs/sql_builder.html

	var transactionRecords []*model.TransactionRecord

	// Method1: Use Raw SQL with scan to query
	repo.db.Debug().Raw(`
	SELECT date, time, stockNo, tranType, quantity, unitPrice
	FROM tblTransactionRecord
	UNION SELECT * FROM tblTransactionRecordSys
	`).Scan(&transactionRecords)

	fmt.Println(transactionRecords[0], "\n\n---")

	var transactionRecords2 []*model.TransactionRecord
	// Method2: Combine GORM API build Raw SQL
	repo.db.Debug().Raw("? UNION ?",
		repo.db.Select("date, time, stockNo, tranType, quantity, unitPrice").Table("tblTransactionRecord"),
		repo.db.Select("*").Table("tblTransactionRecordSys"),
	).Scan(&transactionRecords2)

	fmt.Println(transactionRecords2[0], "\n\n---")

	// ---
	// Sometimes you would like to generate SQL without executing.
	// The Debug function is not what you expect, you can do as below.

	// Method1: Use Raw SQL without scan
	query := repo.db.Raw("SELECT * Error SQL Syntax")

	fmt.Println(query.Statement.SQL.String(), "\n\n---")

	// Method2: Use Raw SQL without scan
	stmt := repo.db.Session(&gorm.Session{DryRun: true}).Select("*").
		Table("Error SQL Syntax").Find(&model.TransactionRecord{}).Statement

	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars, "\n\n---")

	// Method3: Use ToSQL
	sql := repo.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Select("*").Table("Error SQL Syntax").Find(&model.TransactionRecord{})
	})

	fmt.Println(sql)
}

// QueryTransactionRecordUnion
func (repo *TransactionRepositoryGorm) QueryTransactionRecordUnion() ([]*model.TransactionRecord, error) {
	// SQLite seems to help you sort items by primary key when you query via UNION keyword.
	// Or you can add ORDER keyword in the last line to sort it.
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Debug().Raw(`
	SELECT date, time, stockNo, tranType, quantity, unitPrice
	FROM tblTransactionRecord
	UNION SELECT * FROM tblTransactionRecordSys
	`).Scan(&transactionRecords).Error

	if err != nil {
		return nil, nil
	}

	return transactionRecords, nil
}

// insertTransactionRecordSys
func (repo *TransactionRepositoryGorm) InsertTransactionRecordSys(tr *model.TransactionRecord) error {
	if err := repo.db.Debug().Table("tblTransactionRecordSys").Create(tr).Error; err != nil {
		return err
	}

	return nil
}

// Service Tier

// addTransactionTailRecursion add new transaction records with tail recursion,
// When adding, inventory and transaction history, especially write-offs and
// tails, need to be considered.
func (repo *TransactionRepositoryGorm) addTransactionTailRecursion(newTransaction *model.Transaction, remainingQuantity int) (*model.Transaction, error) {
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
		if err != gorm.ErrRecordNotFound {
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
func (repo *TransactionRepositoryGorm) AddTransaction(newTransaction *model.Transaction) (*model.Transaction, error) {
	remainingQuantity := newTransaction.Quantity
	return repo.addTransactionTailRecursion(newTransaction, remainingQuantity)
}
