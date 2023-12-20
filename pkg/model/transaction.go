package model

type Transaction struct {
	ID          int
	StockNo     string
	Date        string
	Quantity    int
	TranType    int
	UnitPrice   float64
	TotalAmount int
	Taxes       int
}

// New Transaction From DB
func NewTransactionFromDB(
	id int, stockNo string, date string, quantity int, tranType int,
	unitPrice float64, totalAmount int, taxes int) *Transaction {
	return &Transaction{
		ID:          id,
		StockNo:     stockNo,
		Date:        date,
		Quantity:    quantity,
		TranType:    tranType,
		UnitPrice:   unitPrice,
		TotalAmount: totalAmount,
		Taxes:       taxes,
	}
}

// New Transaction From User Input
func NewTransactionFromInput(
	stockNo string, date string, quantity int, tranType int, unitPrice float64) *Transaction {
	t := &Transaction{
		StockNo:   stockNo,
		Date:      date,
		Quantity:  quantity,
		TranType:  tranType,
		UnitPrice: unitPrice,
	}
	t.calculateTotalAmount()
	t.calculateTaxesFromTotalAmount()
	return t
}

// Calculate TotalAmount
func (t *Transaction) calculateTotalAmount() {
	t.TotalAmount = int(float64(t.Quantity) * t.UnitPrice)
}

// Calculate Taxes From Quantity And Price
func (t *Transaction) calculateTaxesFromQuantityAndPrice() {
	var taxRate float64 = 0.003
	t.Taxes = int(float64(t.Quantity) * t.UnitPrice * taxRate)
}

// Calculate Taxes From Total Amount
func (t *Transaction) calculateTaxesFromTotalAmount() {
	var taxRate float64 = 0.003
	t.Taxes = int(float64(t.TotalAmount) * taxRate)
}

func (t *Transaction) CalculateTotalAmount() {
	t.calculateTotalAmount()
}
func (t *Transaction) CalculateTaxesFromTotalAmount() {
	t.calculateTaxesFromTotalAmount()
}
