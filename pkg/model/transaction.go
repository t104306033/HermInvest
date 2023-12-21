package model

// Transaction represents a share transaction.
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
	stockNo string, date string, quantity int, tranType int, unitPrice float64) *Transaction {
	t := &Transaction{
		StockNo:   stockNo,
		Date:      date,
		Quantity:  quantity,
		TranType:  tranType,
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

	// Recalculate total amount and taxes
	t.calculateTotalAmount()
	t.calculateTaxes()
}
