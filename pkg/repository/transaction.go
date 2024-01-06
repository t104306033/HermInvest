package repository

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

/******************************************************************************
 *                             Transaction Table                              *
 ******************************************************************************/

// CreateTransaction: insert transaction and return inserted id
func (repo *TransactionRepository) CreateTransaction(t *model.Transaction) (int, error) {
	result := repo.db.Debug().Table("tblTransaction").Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// CreateTransactions: insert transactions and return inserted ids
func (repo *TransactionRepository) CreateTransactions(ts []*model.Transaction) ([]int, error) {
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

func (repo *TransactionRepository) FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Where("stockNo = ?", stockNo).
		Order("date ASC, time ASC").First(&transaction).Error
	if err != nil {
		return &model.Transaction{}, err
	}

	return &transaction, nil
}

// QueryTransactionAll
func (repo *TransactionRepository) QueryTransactionAll() ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionByID
func (repo *TransactionRepository) QueryTransactionByID(id int) (*model.Transaction, error) {
	var transaction *model.Transaction
	err := repo.db.Debug().Table("tblTransaction").Take(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// QueryTransactionByDetails
func (repo *TransactionRepository) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
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
func (repo *TransactionRepository) UpdateTransaction(id int, t *model.Transaction) error {
	err := repo.db.Debug().Table("tblTransaction").Updates(t).Error
	if err != nil {
		return err
	}
	return nil
}

// deleteAlltblTransaction
func (repo *TransactionRepository) DeleteAlltblTransaction() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransaction").Error; err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepository) DeleteTransaction(id int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{ID: id})
	return result.Error
}

func (repo *TransactionRepository) DeleteTransactions(ids []int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{}, "id IN ?", ids)
	return result.Error
}

/******************************************************************************
 *                         Transaction History Table                          *
 ******************************************************************************/

// CreateTransactionHistory: insert transaction and return inserted id
func (repo *TransactionRepository) CreateTransactionHistory(t *model.Transaction) (int, error) {
	result := repo.db.Debug().Table("tblTransactionHistory").Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// CreateTransactionHistorys: insert transactions and return inserted ids
func (repo *TransactionRepository) CreateTransactionHistorys(ts []*model.Transaction) ([]int, error) {
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

// deleteAlltblTransactionHistory
func (repo *TransactionRepository) DeleteAlltblTransactionHistory() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransactionHistory").Error; err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                          Capital Reduction Table                           *
 ******************************************************************************/

// QueryCapitalReductionAll
func (repo *TransactionRepository) QueryCapitalReductionAll() ([]*model.CapitalReduction, error) {
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

/******************************************************************************
 *                          Transaction Record Table                          *
 ******************************************************************************/

// insertTransactionRecordSys
func (repo *TransactionRepository) InsertTransactionRecordSys(tr *model.TransactionRecord) error {
	if err := repo.db.Debug().Table("tblTransactionRecordSys").Create(tr).Error; err != nil {
		return err
	}

	return nil
}

// QueryTransactionRecordByStockNo
func (repo *TransactionRepository) QueryTransactionRecordByStockNo(stockNo string, date string) ([]*model.TransactionRecord, error) {
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

// QueryTransactionRecordUnion
func (repo *TransactionRepository) QueryTransactionRecordUnion() ([]*model.TransactionRecord, error) {
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

// deleteAllTransactionRecordSys
func (repo *TransactionRepository) DeleteAllTransactionRecordSys() error {
	if err := repo.db.Debug().Exec("DELETE FROM tblTransactionRecordSys").Error; err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                                    Note                                    *
 ******************************************************************************/

// QueryUnionNote
func (repo *TransactionRepository) QueryUnionNote() {
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
