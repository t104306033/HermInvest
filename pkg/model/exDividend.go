package model

// Transaction represents a share transaction.
type ExDividend struct {
	YQ               string  `gorm:"column:YQ"`
	StockNo          string  `gorm:"column:stockNo"`
	ExDividendDate   string  `gorm:"column:exDividendDate"`
	DistributionDate string  `gorm:"column:distributionDate"`
	CashDividend     float64 `gorm:"column:cashDividend"`
	StockDividend    float64 `gorm:"column:stockDividend"`
}

func (ed *ExDividend) CalcTransactionRecords(totalQuantity int) *TransactionRecord {
	distributionRecord := ed.calcDistributionRecord(totalQuantity)
	return distributionRecord
}

func (ed *ExDividend) calcDistributionRecord(totalQuantity int) *TransactionRecord {
	return NewTransactionRecord(
		ed.DistributionDate, "08:00:10",
		ed.StockNo, -1, totalQuantity, ed.CashDividend)
}
