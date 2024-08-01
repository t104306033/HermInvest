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

func (cr *CapitalReduction) TableName() string {
	return "tblCapitalReduction" // default table name
}

func (cr *CapitalReduction) CalcTransactionRecords(totalQuantity int, avgUnitPrice float64) (*TransactionRecord, *TransactionRecord) {
	capitalReductionRecord := cr.calcCapitalReductionRecord(totalQuantity, avgUnitPrice)
	distributionRecord := cr.calcDistributionRecord(totalQuantity, avgUnitPrice)
	return capitalReductionRecord, distributionRecord
}

func (cr *CapitalReduction) calcCapitalReductionRecord(totalQuantity int, avgUnitPrice float64) *TransactionRecord {
	return NewTransactionRecord(
		cr.CapitalReductionDate, "08:00:00",
		cr.StockNo, -1, totalQuantity, avgUnitPrice)
}

func (cr *CapitalReduction) calcDistributionRecord(totalQuantity int, avgUnitPrice float64) *TransactionRecord {
	if cr.NewStockNo == "" {
		cr.NewStockNo = cr.StockNo
	}

	distributionQuantity := int(float64(totalQuantity) * (1 - cr.Ratio))
	distributionUnitPrice := (avgUnitPrice - cr.Cash) / (1 - cr.Ratio)

	return NewTransactionRecord(
		cr.DistributionDate, "08:00:10",
		cr.NewStockNo, 1, distributionQuantity, distributionUnitPrice)
}
