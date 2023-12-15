package main

type Transaction struct {
	id          int
	stockNo     string
	date        string
	quantity    int
	tranType    int
	unitPrice   float64
	totalAmount int
	taxes       int
}

// New Transaction From DB
func newTransactionFromDB(
	id int, stockNo string, date string, quantity int, tranType int,
	unitPrice float64, totalAmount int, taxes int) *Transaction {
	return &Transaction{
		id:          id,
		stockNo:     stockNo,
		date:        date,
		quantity:    quantity,
		tranType:    tranType,
		unitPrice:   unitPrice,
		totalAmount: totalAmount,
		taxes:       taxes,
	}
}

// New Transaction From User Input
func newTransactionFromInput(
	stockNo string, date string, quantity int, tranType int, unitPrice float64) *Transaction {
	t := &Transaction{
		stockNo:   stockNo,
		date:      date,
		quantity:  quantity,
		tranType:  tranType,
		unitPrice: unitPrice,
	}
	t.calculateTotalAmount()
	t.calculateTaxesFromTotalAmount()
	return t
}

// Calculate TotalAmount
func (t *Transaction) calculateTotalAmount() {
	t.totalAmount = int(float64(t.quantity) * t.unitPrice)
}

// Calculate Taxes From Quantity And Price
func (t *Transaction) calculateTaxesFromQuantityAndPrice() {
	var taxRate float64 = 0.003
	t.taxes = int(float64(t.quantity) * t.unitPrice * taxRate)
}

// Calculate Taxes From Total Amount
func (t *Transaction) calculateTaxesFromTotalAmount() {
	var taxRate float64 = 0.003
	t.taxes = int(float64(t.totalAmount) * taxRate)
}
