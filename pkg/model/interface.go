package model

import (
	"gorm.io/gorm"
)

type Repositorier interface {
	CreateTransaction(t *Transaction) (int, error)
	CreateTransactionHistory(t *Transaction) (int, error)
	CreateTransactionHistorys(ts []*Transaction) ([]int, error)
	CreateTransactions(ts []*Transaction) ([]int, error)
	DeleteAllTransactionRecordSys() error
	DeleteAlltblTransaction() error
	DeleteAlltblTransactionHistory() error
	DeleteTransaction(id int) error
	DeleteTransactions(ids []int) error
	FindEarliestTransactionByStockNo(stockNo string) (*Transaction, error)
	InsertTransactionRecordSys(tr *TransactionRecord) error
	QueryCapitalReductionAll() ([]*CapitalReduction, error)
	QueryTransactionAll() ([]*Transaction, error)
	QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*Transaction, error)
	QueryTransactionByID(id int) (*Transaction, error)
	QueryTransactionRecordByStockNo(stockNo string, date string) ([]*TransactionRecord, error)
	QueryTransactionRecordUnion() ([]*TransactionRecord, error)
	QueryUnionNote()
	UpdateTransaction(id int, t *Transaction) error

	WithTrx(trxHandle *gorm.DB) Repositorier
	Begin() *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
}
