package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/repository"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. calc total amount and taxes
// 3. build sql syntax
// 4. insert into sql
// 5. print out result

var addCmd = &cobra.Command{
	Use:   "add stockNo type quantity unitPrice [date]",
	Short: "Add stock (Stock No., Type, Quantity, Unit Price)",
	Example: "" +
		"  - Purchase at today's date:\n" +
		"    hermInvestCli stock add 0050 1 1500 23.5\n\n" +

		"  - Sale on a specific date:\n" +
		"    hermInvestCli stock add -- 0050 -1 1500 23.5 2023-12-01",
	Long: `Add stock by transaction stock number, type, quantity, and unit price`,
	Args: cobra.RangeArgs(4, 5),
	Run:  addRun,
}

func init() {
	stockCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	stockNo, tranType, quantity, unitPrice, date, err := ParseTransactionForAddCmd(args)
	if err != nil {
		fmt.Println("Error parsing transaction data:", err)
		return
	}

	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}
	defer db.Close()

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	// add stock in inventory
	// 1. new transaction from input
	// 2. find the first purchase from the inventory
	// 3. check the transaction type of new transaction and first purchase

	// TODO: service.addTransaction() AddTransactionAndUpdateInventory
	newTransaction := model.NewTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)

	// Find the first purchase from the inventory
	earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
	if err == sql.ErrNoRows {
		// If no first purchase found, set TranType and handle accordingly
		earliestTransaction.TranType = newTransaction.TranType
	} else if err != nil {
		fmt.Println("Error finding first purchase:", err)
		return
	}

	if earliestTransaction.TranType == newTransaction.TranType {
		// add stock in the inventory directly
		id, err := repo.CreateTransaction(newTransaction)
		if err != nil {
			fmt.Println("Error creating transaction: ", err)
			return
		}
		transactions, err := repo.QueryTransactionByID(id)
		if err != nil {
			fmt.Println("Error querying database:", err)
			return
		}

		// Print out result
		displayResults(transactions)

		return
	}

	inventoryTransactions, err := repo.QueryInventoryTransactions(newTransaction.StockNo, newTransaction.Quantity)
	if err != nil {
		// handle query old purchase failed
		fmt.Println("Error querying old purchase:", err)
		return
	}

	remainingQuantity := newTransaction.Quantity
	for _, t := range inventoryTransactions {
		remainingQuantity -= t.Quantity
	}

	err = repo.MoveInventoryToTransactionHistorys(inventoryTransactions)
	if err != nil {
		fmt.Println("Error moving inventory to transaction history:", err)
		return
	}

	if remainingQuantity > 0 {
		// Find the first purchase again after old purchases are processed
		earliestTransaction, err := repo.FindEarliestTransactionByStockNo(newTransaction.StockNo)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("add stock in the inventory directly", err)

				// record newTransaction to history
				newTransaction.SetQuantity(newTransaction.Quantity - remainingQuantity)
				_, err = repo.CreateTransactionHistory(newTransaction)
				if err != nil {
					// handle create transaction history failed
				}

				newTransaction.SetQuantity(remainingQuantity)
				id, err := repo.CreateTransaction(newTransaction)
				if err != nil {
					fmt.Println("Error creating transaction: ", err)
					return
				}

				transactions, err := repo.QueryTransactionByID(id)
				if err != nil {
					fmt.Println("Error querying database:", err)
					return
				}
				// Print out result
				displayResults(transactions)

				fmt.Println(newTransaction.Quantity, remainingQuantity)

			} else {
				// Pass, can't find first purchase, cause no stock in the inventory
				fmt.Println("Error finding earliest transaction:", err)
			}
			return
		}

		if earliestTransaction.Quantity < remainingQuantity {
			fmt.Println("Not Impelment yet: earliestTransaction.Quantity < remainingQuantity: ", err)
			return
		}

		// Update the quantity of the first purchase
		earliestTransaction.SetQuantity(earliestTransaction.Quantity - remainingQuantity)
		err = repo.UpdateTransaction(earliestTransaction.ID, earliestTransaction)
		if err != nil {
			// handle update transaction failed
		}

		// record purchase history
		earliestTransaction.SetQuantity(remainingQuantity)
		_, err = repo.CreateTransactionHistory(earliestTransaction)
		if err != nil {
			// handle create transaction history failed
		}
	}

	// record newTransaction to history
	_, err = repo.CreateTransactionHistory(newTransaction)
	if err != nil {
		// handle create transaction history failed
	}

}

func ParseTransactionForAddCmd(args []string) (string, int, int, float64, string, error) {
	stockNo := args[0] // regex a-z 0-9

	tranType, err := strconv.Atoi(args[1])
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing integer: %s", err)
	}

	quantity, err := strconv.Atoi(args[2])
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing integer: %s", err)
	}

	unitPrice, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return "", 0, 0, 0, "", fmt.Errorf("error parsing float: %s", err)
	}

	var date string
	if len(args) > 4 {
		parsedTime, err := time.Parse(time.DateOnly, args[4])
		if err != nil {
			return "", 0, 0, 0, "", fmt.Errorf("error parsing date: %s", err)
		}
		date = parsedTime.Format(time.DateOnly)
	} else {
		date = time.Now().Format(time.DateOnly)
	}

	return stockNo, tranType, quantity, unitPrice, date, nil
}
