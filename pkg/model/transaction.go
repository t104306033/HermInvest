package model

// Transaction represents a share transaction.
type Transaction struct {
	ID          int
	Date        string
	Time        string
	StockNo     string
	TranType    int
	Quantity    int
	UnitPrice   float64
	TotalAmount int
	Taxes       int
}

// NewTransactionFromDB creates a new Transaction object from database records.
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

// NewTransactionFromInput creates a new Transaction object from input.
// It initializes the transaction with inputs. Additionally, the total amount
// and taxes are recalculated based on the new transaction details.
func NewTransactionFromInput(
	date string, time string, stockNo string, tranType int, quantity int,
	unitPrice float64) *Transaction {
	t := &Transaction{
		Date:      date,
		Time:      time,
		StockNo:   stockNo,
		TranType:  tranType,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
	t.calculateTotalAmount()
	t.calculateTaxes()
	return t
}

// calculateTotalAmount calculates the total amount based on transaction details.
func (t *Transaction) calculateTotalAmount() {
	t.TotalAmount = int(float64(t.Quantity) * t.UnitPrice)
}

// calculateTotalAmount calculates the taxes based on transaction details.
func (t *Transaction) calculateTaxes() {
	var taxRate float64 = 0.003
	t.Taxes = int(float64(t.TotalAmount) * taxRate)
}

// SetUnitPrice updates the unit price of the transaction.
// It recalculates the total amount and taxes based on the updated unit price.
// The calculation of total amount and taxes are interdependent.
func (t *Transaction) SetUnitPrice(unitPrice float64) {
	t.UnitPrice = unitPrice

	t.recalculate()
}

// SetQuantity updates the quantity of the transaction.
// It recalculates the total amount and taxes based on the quantity.
// The calculation of total amount and taxes are interdependent.
func (t *Transaction) SetQuantity(quantity int) {
	t.Quantity = quantity

	t.recalculate()
}

// recalculate total amount and taxes of the transaction.
// It will recalculates the total amount and taxes based on the model.
func (t *Transaction) recalculate() {
	t.calculateTotalAmount()
	t.calculateTaxes()
}
