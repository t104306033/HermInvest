package main

import (
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
		"    hermInvestCli stock add 0050 -1 1500 23.5 2023-12-01",
	Long: `Add stock by transaction stock number, type, quantity, and unit price`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		stockNo := args[0] // regex a-z 0-9
		tranType, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
			return
		}
		quantity, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
			return
		}
		unitPrice, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Println("Error parsing float: ", err)
			return
		}

		var date string
		if len(args) > 4 {
			// check user input time format is correct
			parsedTime, err := time.Parse(time.DateOnly, args[4])
			if err != nil {
				fmt.Println("Error parsing date: ", err)
				return
			}
			date = parsedTime.Format(time.DateOnly)
		} else {
			date = time.Now().Format(time.DateOnly)
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error geting DB connection: ", err)
		}
		defer db.Close()

		// init transactionRepository
		repo := &transactionRepository{db: db}
		t := newTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)
		id, err := repo.createTransaction(t)
		if err != nil {
			fmt.Println("Error creating transaction: ", err)
		}

		transactions, err := repo.queryTransactionByID(id)
		if err != nil {
			fmt.Println("Error querying database:", err)
		}

		// Print out result
		displayResults(transactions)

	},
}
