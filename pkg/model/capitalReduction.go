package model

// Transaction represents a share transaction.
type CapitalReduction struct {
	YQ                   string
	StockNo              string
	CapitalReductionDate string
	DistributionDate     string
	Cash                 float64
	Ratio                float64
	NewStockNo           string
}
