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
	QueryTransactionByID(id int) (*Transaction, error)
	QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error)
	QueryTransactionRecordAll() ([]*TransactionRecord, error)
	QueryTransactionRecordSysAll() ([]*TransactionRecord, error)

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
