package service

import "HermInvest/pkg/model"

type repositorier interface {
	CreateTransaction(t *model.Transaction) (int, error)
	CreateTransactionHistory(t *model.Transaction) (int, error)
	CreateTransactionHistorys(ts []*model.Transaction) ([]int, error)
	CreateTransactions(ts []*model.Transaction) ([]int, error)
	DeleteAllTransactionRecordSys() error
	DeleteAlltblTransaction() error
	DeleteAlltblTransactionHistory() error
	DeleteTransaction(id int) error
	DeleteTransactions(ids []int) error
	FindEarliestTransactionByStockNo(stockNo string) (*model.Transaction, error)
	InsertTransactionRecordSys(tr *model.TransactionRecord) error
	QueryCapitalReductionAll() ([]*model.CapitalReduction, error)
	QueryTransactionAll() ([]*model.Transaction, error)
	QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error)
	QueryTransactionByID(id int) (*model.Transaction, error)
	QueryTransactionRecordByStockNo(stockNo string, date string) ([]*model.TransactionRecord, error)
	QueryTransactionRecordUnion() ([]*model.TransactionRecord, error)
	QueryUnionNote()
	UpdateTransaction(id int, t *model.Transaction) error
}
