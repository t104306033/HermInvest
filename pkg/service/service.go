package service

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

type service struct {
	repo model.Repositorier
}

func NewService(repository model.Repositorier) *service {
	return &service{repo: repository}
}

// addTransactionTailRecursion add new transaction records with tail recursion,
// When adding, inventory and transaction history, especially write-offs and
// tails, need to be considered.
func (serv *service) addTransactionTailRecursion(newTransaction *model.Transaction, remainingQuantity int) (*model.Transaction, error) {
	// Principles:
	// 1. Ensure that each transaction has a corresponding transaction record.
	// 2. Update inventory quantities based on transactions, including adding,
	//    reducing, or deleting inventory.
	// 3. Depending on the transaction situation, only transaction history can
	//    be added and cannot be modified or deleted.
	// 4. For insufficient write-off quantities, recursive processing is used
	//    to ensure that the write-off is completed.

	// Cases:
	// 1. Newly added: If there is no transaction in the inventory (A) or
	//    the new transaction is the same as the oldest transaction in the
	//    inventory (B), add it directly to the inventory.
	// 2. Write-off:
	// 	* Sufficient inventory: If the inventory quantity is sufficient,
	//    update the inventory quantity (C) or delete the inventory (D), and
	//    add the corresponding transaction history.
	// 	* Insufficient inventory: If the inventory quantity can't be Write-off.
	//    Recurse until success (E). The termination condition is A B C D.
	//  * Over inventory: Write-off over than inventory (F).

	// TODO: This func should be moved to service tier.

	earliestTransaction, err := serv.repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("error finding first purchase: %v", err)
		}
		// Case A
		earliestTransaction.TranType = newTransaction.TranType
	}

	if earliestTransaction.TranType == newTransaction.TranType {
		if newTransaction.Quantity != remainingQuantity {
			// Case F
			newTransaction.SetQuantity(newTransaction.Quantity - remainingQuantity)
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			newTransaction.SetQuantity(remainingQuantity)
		}

		// Case B
		id, err := serv.repo.CreateTransaction(newTransaction)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %v", err)
		}
		transaction, err := serv.repo.QueryTransactionByID(id)
		if err != nil {
			return nil, fmt.Errorf("error querying database: %v", err)
		}

		return transaction, nil
	} else {
		if earliestTransaction.Quantity > remainingQuantity {
			// Case C

			// Create a copy for adding stock history
			stockHistoryAdd := &model.Transaction{}
			*stockHistoryAdd = *earliestTransaction
			// var stockHistoryAdd *model.Transaction // why can't use it, study it
			// *stockHistoryAdd = *earliestTransaction

			// add transaction history
			stockHistoryAdd.SetQuantity(remainingQuantity)
			_, err = serv.repo.CreateTransactionHistory(stockHistoryAdd)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// Update stock inventory
			earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
			err := serv.repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
			if err != nil {
				// handle update transaction failed
			}

			return earliestTransaction, nil
		} else if earliestTransaction.Quantity == remainingQuantity {
			// Case D

			// add transaction history
			_, err = serv.repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = serv.repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			// delete stock inventory
			err = serv.repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			// Or use move

			return nil, nil
		} else { // earliestTransaction.Quantity < remainingQuantity
			// Case E

			// add transaction history
			_, err = serv.repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// delete stock inventory
			err = serv.repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			remainingQuantity = remainingQuantity - earliestTransaction.Quantity

			return serv.addTransactionTailRecursion(newTransaction, remainingQuantity)
		}
	}
}

// AddTransaction add the transaction from the input to the inventory.
// It will add or update transactions in the inventory and add history.
// Return the modified transaction record in the inventory
func (serv *service) AddTransaction(newTransaction *model.Transaction) (*model.Transaction, error) {
	remainingQuantity := newTransaction.Quantity
	return serv.addTransactionTailRecursion(newTransaction, remainingQuantity)
}

// ---

func (serv *service) DeleteTransaction(id int) error {
	return serv.repo.DeleteTransaction(id)
}

func (serv *service) QueryTransactionAll() ([]*model.Transaction, error) {
	return serv.repo.QueryTransactionAll()
}

func (serv *service) QueryTransactionByID(id int) (*model.Transaction, error) {
	return serv.repo.QueryTransactionByID(id)
}

func (serv *service) QueryTransactionByDetails(stockNo string, tranType int, date string) ([]*model.Transaction, error) {
	return serv.repo.QueryTransactionByDetails(stockNo, tranType, date)
}

func (serv *service) UpdateTransaction(id int, t *model.Transaction) error {
	return serv.repo.UpdateTransaction(id, t)
}

func (serv *service) DeleteAllTransactionRecordSys() error {
	return serv.repo.DeleteAllTransactionRecordSys()
}

func (serv *service) InsertTransactionRecordSys(tr *model.TransactionRecord) error {
	return serv.repo.InsertTransactionRecordSys(tr)
}

func (serv *service) DeleteAlltblTransaction() error {
	return serv.repo.DeleteAlltblTransaction()
}

func (serv *service) DeleteAlltblTransactionHistory() error {
	return serv.repo.DeleteAlltblTransactionHistory()
}

func (serv *service) QueryTransactionRecordUnion() ([]*model.TransactionRecord, error) {
	return serv.repo.QueryTransactionRecordUnion()
}
func (serv *service) CreateTransactions(ts []*model.Transaction) ([]int, error) {
	return serv.repo.CreateTransactions(ts)
}

// ---

func (serv *service) RebuildCapitalReduction() error {

	tx := serv.repo.Begin()

	serv.repo.WithTrx(tx).DeleteAllTransactionRecordSys()

	// 1. Query all transaction records from tblCapitalReduction
	crs, err := serv.repo.WithTrx(tx).QueryCapitalReductionAll()
	if err != nil {
		serv.repo.WithTrx(tx).Rollback()
		return err
	}

	// 2. Iterate over each capital reduction record
	for _, cr := range crs {
		// Query transaction records by stock number
		trs, err := serv.repo.WithTrx(tx).QueryTransactionRecordByStockNo(cr.StockNo, cr.CapitalReductionDate)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return err
		}

		remainingTrs := make([]*model.TransactionRecord, 0)
		for _, tr := range trs {
			if tr.TranType > 0 {
				remainingTrs = append(remainingTrs, tr)
			} else {
				qty := tr.Quantity
				for qty > 0 && len(remainingTrs) > 0 {
					var remove *model.TransactionRecord
					remove, remainingTrs = remainingTrs[0], remainingTrs[1:]
					qty -= remove.Quantity
				}
			}
		}

		var totalQuantity int
		var totalAmount int
		for _, tr := range remainingTrs {
			totalQuantity += tr.TranType * tr.Quantity
			totalAmount += int(float64(tr.Quantity) * tr.UnitPrice)
		}
		var avgUnitPrice float64 = float64(totalAmount) / float64(totalQuantity)

		// 3. insert into tblTransactionRecordSys
		capitalReductionRecord := model.NewTransactionRecord(
			cr.CapitalReductionDate, "08:00:00", cr.StockNo, -1, totalQuantity, avgUnitPrice)
		err = serv.repo.WithTrx(tx).InsertTransactionRecordSys(capitalReductionRecord)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return err
		}

		newStockNo := cr.NewStockNo
		if newStockNo == "" {
			newStockNo = cr.StockNo
		}

		newQuantity := int(float64(totalQuantity) * (1 - cr.Ratio))
		newAvgUnitPrice := (avgUnitPrice - cr.Cash) / (1 - cr.Ratio)
		newCapitalReductionRecord := model.NewTransactionRecord(
			cr.DistributionDate, "08:00:10", newStockNo, 1, newQuantity, newAvgUnitPrice)
		err = serv.repo.WithTrx(tx).InsertTransactionRecordSys(newCapitalReductionRecord)
		if err != nil {
			serv.repo.WithTrx(tx).Rollback()
			return err
		}
	}

	serv.repo.WithTrx(tx).Commit()
	return nil

}
