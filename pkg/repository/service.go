package repository

import (
	"HermInvest/pkg/model"
	"fmt"

	"gorm.io/gorm"
)

// Service Tier

// addTransactionTailRecursion add new transaction records with tail recursion,
// When adding, inventory and transaction history, especially write-offs and
// tails, need to be considered.
func (repo *TransactionRepository) addTransactionTailRecursion(newTransaction *model.Transaction, remainingQuantity int) (*model.Transaction, error) {
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

	earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
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
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			newTransaction.SetQuantity(remainingQuantity)
		}

		// Case B
		id, err := repo.CreateTransaction(newTransaction)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %v", err)
		}
		transaction, err := repo.QueryTransactionByID(id)
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
			_, err = repo.CreateTransactionHistory(stockHistoryAdd)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// Update stock inventory
			earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
			err := repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
			if err != nil {
				// handle update transaction failed
			}

			return earliestTransaction, nil
		} else if earliestTransaction.Quantity == remainingQuantity {
			// Case D

			// add transaction history
			_, err = repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			_, err = repo.CreateTransactionHistory(newTransaction)
			if err != nil {
				// handle create transaction history failed
			}
			// delete stock inventory
			err = repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			// Or use move

			return nil, nil
		} else { // earliestTransaction.Quantity < remainingQuantity
			// Case E

			// add transaction history
			_, err = repo.CreateTransactionHistory(earliestTransaction)
			if err != nil {
				// handle create transaction history failed
			}

			// delete stock inventory
			err = repo.DeleteTransaction(earliestTransaction.ID)
			if err != nil {
				// handle create transaction history failed
			}

			remainingQuantity = remainingQuantity - earliestTransaction.Quantity

			return repo.addTransactionTailRecursion(newTransaction, remainingQuantity)
		}
	}
}

// AddTransaction add the transaction from the input to the inventory.
// It will add or update transactions in the inventory and add history.
// Return the modified transaction record in the inventory
func (repo *TransactionRepository) AddTransaction(newTransaction *model.Transaction) (*model.Transaction, error) {
	remainingQuantity := newTransaction.Quantity
	return repo.addTransactionTailRecursion(newTransaction, remainingQuantity)
}
