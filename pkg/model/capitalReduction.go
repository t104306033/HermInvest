package model

// Transaction represents a share transaction.
type CapitalReduction struct {
	YQ                   string  `gorm:"column:YQ"`
	StockNo              string  `gorm:"column:stockNo"`
	CapitalReductionDate string  `gorm:"column:capitalReductionDate"`
	DistributionDate     string  `gorm:"column:distributionDate"`
	Cash                 float64 `gorm:"column:cash"`
	Ratio                float64 `gorm:"column:ratio"`
	NewStockNo           string  `gorm:"column:newStockNo"`
}
