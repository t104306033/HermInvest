package model

import (
	"gorm.io/gorm"
)

type Repositorier interface {
	CreateTransaction(t *Transaction) (int, error)
	CreateTransactionHistory(t *Transaction) (int, error)
	// CreateTransactionHistorys(ts []*Transaction) ([]int, error)
	CreateTransactions(ts []*Transaction) ([]int, error)
	// DeleteAllTransactionRecordSys() error
	// DeleteAlltblTransaction() error
	// DeleteAlltblTransactionHistory() error
	// DeleteAllCashDividendRecord() error
	DeleteTransaction(id int) error
	// DeleteSQLiteSequence() error
	DeleteTransactions(ids []int) error
	FindEarliestTransactionByStockNo(stockNo string) (*Transaction, error)
	InsertTransactionRecordSys(tr *TransactionRecord) error
	InsertCashDividendRecord(cd *ExDividend) error
	QueryCapitalReductionAll() ([]*CapitalReduction, error)
	QueryDividendAll() ([]*ExDividend, error)
	QueryTransactionAll() ([]*Transaction, error)
	QueryTransactionRecordSysAll() ([]*TransactionRecord, error)
	QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error)
	QueryTransactionByID(id int) (*Transaction, error)
	QueryTransactionRecordAll() ([]*TransactionRecord, error)
	QueryTransactionRecordByStockNo(stockNo string, date string) ([]*TransactionRecord, error)
	QueryTransactionRecordUnion() ([]*TransactionRecord, error)
	QueryUnionNote()
	UpdateTransaction(id int, t *Transaction) error
	WithTrx(trxHandle *gorm.DB) Repositorier
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB

	// tablename: "sqlite_sequence", "tblTransaction", "tblTransactionHistory", "tblTransactionCash", "tblTransactionRecordSys"
	DropTable(tablename string) error
}
