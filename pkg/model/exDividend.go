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
