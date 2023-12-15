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
		"    hermInvestCli stock add 0050 -1 1500 23.5 2023-12-01\n\n" +

		"  - Purchase on a specific date assign id:\n" +
		"    hermInvestCli stock add --id 33 0050 1 1500 23.5 2023-12-01",
	Long: `Add stock by transaction stock number, type, quantity, and unit price`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
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

		t := NewTransactionFromUserInput(id, stockNo, date, quantity, tranType, unitPrice)

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error geting DB connection: ", err)
		}
		defer db.Close()

		// Execute the insert query
		if id != 0 {
			query := `INSERT INTO tblTransaction (id, stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
			_, err = db.Exec(query, t.id, t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Pass: Stock added successfully!")
			}

			// Print out result
			rows, err := db.Query(buildQueryByID(id))
			if err != nil {
				fmt.Println("Error querying database:", err)
				return
			}
			defer rows.Close()

			displayResults(rows)
		} else {
			query := `INSERT INTO tblTransaction (stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?)`
			_, err = db.Exec(query, t.stockNo, t.date, t.quantity, t.tranType, t.unitPrice, t.totalAmount, t.taxes)
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println("Pass: Stock added successfully!")
			}
		}
	},
}

func init() {
	addCmd.Flags().Int("id", 0, "Add by ID")
}
