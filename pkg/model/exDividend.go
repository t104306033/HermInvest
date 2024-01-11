package model

// Transaction represents a share transaction.
type ExDividend struct {
	YQ               string  `gorm:"column:YQ"`
	StockNo          string  `gorm:"column:stockNo"`
	ExDividendDate   string  `gorm:"column:exDividendDate"`
	DistributionDate string  `gorm:"column:distributionDate"`
	CashDividend     float64 `gorm:"column:cashDividend"`
	StockDividend    float64 `gorm:"column:stockDividend"`
	Quantity         int     `gorm:"column:quantity"`
	TotalAmount      int     `gorm:"column:totalAmount"`
}

// NewTransactionRecord creates a new transaction record object.
func NewCashDividendRecord(yq, stockNo, exDividendDate, distributionDate string,
	cashDividend float64, quantity, totalAmount int) *ExDividend {
	return &ExDividend{
		YQ:               yq,
		StockNo:          stockNo,
		ExDividendDate:   exDividendDate,
		DistributionDate: distributionDate,
		CashDividend:     cashDividend,
		Quantity:         quantity,
		TotalAmount:      totalAmount,
	}
}

func (ed *ExDividend) TableName() string {
	return "tblDividend" // default table name
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

func (ed *ExDividend) CalcCashDividendRecord(totalQuantity int) *ExDividend {
	totalAmount := int(float64(totalQuantity) * ed.CashDividend)

	return NewCashDividendRecord(
		ed.YQ, ed.StockNo, ed.ExDividendDate, ed.DistributionDate,
		ed.CashDividend, totalQuantity, totalAmount)
}
