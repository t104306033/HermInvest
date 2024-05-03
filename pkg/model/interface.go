package model

import (
	"gorm.io/gorm"
)

type Repositorier interface {
	CreateTransaction(t *Transaction) (int, error)
	CreateTransactionHistory(t *Transaction) (int, error)
	CreateTransactions(ts []*Transaction) ([]int, error)
	CreateTransactionRecordSys(tr *TransactionRecord) error
	CreateCashDividendRecord(cd *ExDividend) error

	FindEarliestTransactionByStockNo(stockNo string) (*Transaction, error)
	QueryCapitalReductionAll() ([]*CapitalReduction, error)
	QueryDividendAll() ([]*ExDividend, error)
	QueryTransactionAll() ([]*Transaction, error)
	QueryTransactionRecordSysAll() ([]*TransactionRecord, error)
	QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error)
	QueryTransactionByID(id int) (*Transaction, error)
	QueryTransactionRecordAll() ([]*TransactionRecord, error)
	QueryTransactionRecordByStockNo(stockNo string, date string) ([]*TransactionRecord, error)

	UpdateTransaction(id int, t *Transaction) error

	DeleteTransaction(id int) error
	DeleteTransactions(ids []int) error

	DropTable(tablename string) error

	gormRepositorier
}

type gormRepositorier interface {
	WithTrx(trxHandle *gorm.DB) Repositorier
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
