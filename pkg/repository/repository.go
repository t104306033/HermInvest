package repository

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) WithTrx(trxHandle *gorm.DB) model.Repositorier {
	// if trxHandle == nil {
	// 	fmt.Println("WithTrx: Transaction Database not found")
	// 	return repo
	// }
	// fmt.Println("WithTrx: Transaction Database found")
	return &repository{db: trxHandle} // return new one
}

func (repo *repository) Begin() *gorm.DB {
	return repo.db.Begin()
}

func (repo *repository) Commit() *gorm.DB {
	return repo.db.Commit()
}

func (repo *repository) Rollback() *gorm.DB {
	return repo.db.Rollback()
}

// echo "Transaction Table" | boxes -a c -s 80 -d cc

/******************************************************************************
 *                             Transaction Table                              *
 ******************************************************************************/

// CreateTransaction: insert transaction and return inserted id
func (repo *repository) CreateTransaction(t *model.Transaction) (int, error) {
	result := repo.db.Table("tblTransaction").Create(&t)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// CreateTransactions: insert transactions and return inserted ids
func (repo *repository) CreateTransactions(ts []*model.Transaction) ([]int, error) {
	result := repo.db.Table("tblTransaction").Create(&ts)
	if result.Error != nil {
		return nil, result.Error
	}

	var insertedIDs []int
	for _, t := range ts {
		insertedIDs = append(insertedIDs, t.ID)
	}

	return insertedIDs, nil
}

func (repo *repository) FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := repo.db.Table("tblTransaction").Where("stockNo = ?", stockNo).
		Order("date ASC, time ASC").First(&transaction).Error
	if err != nil {
		return &model.Transaction{}, err
	}

	return &transaction, nil
}

// QueryTransactionAll
func (repo *repository) QueryTransactionAll() ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	err := repo.db.Table("tblTransaction").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// QueryTransactionByID
func (repo *repository) QueryTransactionByID(id int) (*model.Transaction, error) {
	var transaction *model.Transaction
	err := repo.db.Table("tblTransaction").Take(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// QueryTransactionByDetails
func (repo *repository) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
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

	err := repo.db.Table("tblTransaction").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// updateTransaction
func (repo *repository) UpdateTransaction(id int, t *model.Transaction) error {
	err := repo.db.Table("tblTransaction").Updates(t).Error
	if err != nil {
		return err
	}
	return nil
}

// deleteAlltblTransaction
func (repo *repository) DeleteAlltblTransaction() error {
	if err := repo.db.Exec("DELETE FROM tblTransaction").Error; err != nil {
		return err
	}

	return nil
}

func (repo *repository) DeleteTransaction(id int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{ID: id})
	return result.Error
}

func (repo *repository) DeleteTransactions(ids []int) error {
	result := repo.db.Table("tblTransaction").Delete(&model.Transaction{}, "id IN ?", ids)
	return result.Error
}

/******************************************************************************
 *                         Transaction History Table                          *
 ******************************************************************************/

// CreateTransactionHistory: insert transaction and return inserted id
func (repo *repository) CreateTransactionHistory(t *model.Transaction) (int, error) {
	// Create a new one and set ID to 0, let SQLite autoincrement
	var transactionHistory model.Transaction = *t
	transactionHistory.ID = 0

	result := repo.db.Table("tblTransactionHistory").Create(&transactionHistory)
	if result.Error != nil {
		return 0, result.Error
	}

	return t.ID, nil
}

// // CreateTransactionHistorys: insert transactions and return inserted ids
// func (repo *repository) CreateTransactionHistorys(ts []*model.Transaction) ([]int, error) {
// 	result := repo.db.Table("tblTransactionHistory").Create(&ts)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	var insertedIDs []int
// 	for _, t := range ts {
// 		insertedIDs = append(insertedIDs, t.ID)
// 	}

// 	return insertedIDs, nil
// }

// deleteAlltblTransactionHistory
func (repo *repository) DeleteAlltblTransactionHistory() error {
	if err := repo.db.Exec("DELETE FROM tblTransactionHistory").Error; err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                                   SQLite                                   *
 ******************************************************************************/

// deleteSQLiteSequence
func (repo *repository) DeleteSQLiteSequence() error {
	if err := repo.db.Exec("DELETE FROM sqlite_sequence").Error; err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                          Capital Reduction Table                           *
 ******************************************************************************/

// QueryCapitalReductionAll
func (repo *repository) QueryCapitalReductionAll() ([]*model.CapitalReduction, error) {
	var capitalReductions []*model.CapitalReduction
	// 使用 Gorm 框架的 Find 方法來執行查詢
	if err := repo.db.Table("tblCapitalReduction").Find(&capitalReductions).Error; err != nil {
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
func (repo *repository) InsertTransactionRecordSys(tr *model.TransactionRecord) error {
	if err := repo.db.Table("tblTransactionRecordSys").Create(tr).Error; err != nil {
		return err
	}

	return nil
}

// QueryTransactionRecordByStockNo
func (repo *repository) QueryTransactionRecordByStockNo(stockNo string, date string) ([]*model.TransactionRecord, error) {
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Table("tblTransactionRecord").Where("stockNo = ? and date < ?", stockNo, date).Find(&transactionRecords).Error
	if err != nil {
		return nil, err
	}

	// for _, cr := range transactionRecords {
	// 	fmt.Println(cr)
	// }

	return transactionRecords, nil
}

// QueryTransactionRecordUnion
func (repo *repository) QueryTransactionRecordUnion() ([]*model.TransactionRecord, error) {
	// SQLite seems to help you sort items by primary key when you query via UNION keyword.
	// Or you can add ORDER keyword in the last line to sort it.
	var transactionRecords []*model.TransactionRecord
	err := repo.db.Raw(`
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
func (repo *repository) DeleteAllTransactionRecordSys() error {
	if err := repo.db.Exec("DELETE FROM tblTransactionRecordSys").Error; err != nil {
		return err
	}

	return nil
}

/******************************************************************************
 *                                    Note                                    *
 ******************************************************************************/

// QueryUnionNote
func (repo *repository) QueryUnionNote() {
	// Most ORMs seem not support UNION keyword, due to its complexity.
	// Faced with this situation, community suggest using "Raw" method to do this.
	// > Reference:
	// > * https://github.com/go-gorm/gorm/issues/3781
	// > * https://stackoverflow.com/questions/67190972/how-to-use-mysql-union-all-on-gorm
	// > * https://gorm.io/docs/sql_builder.html

	var transactionRecords []*model.TransactionRecord

	// Method1: Use Raw SQL with scan to query
	repo.db.Raw(`
	SELECT date, time, stockNo, tranType, quantity, unitPrice
	FROM tblTransactionRecord
	UNION SELECT * FROM tblTransactionRecordSys
	`).Scan(&transactionRecords)

	fmt.Println(transactionRecords[0], "\n\n---")

	var transactionRecords2 []*model.TransactionRecord
	// Method2: Combine GORM API build Raw SQL
	repo.db.Raw("? UNION ?",
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
