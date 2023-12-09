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
	Use:   "add id stockNo type quantity unitPrice [date]",
	Short: "Add stock",
	Long:  `Add stock to the inventory`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		stockNo := args[1]
		tranType := args[2]
		quantity, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
		}
		unitPrice, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Println("Error parsing float: ", err)
		}

		// Parse date argument or default is today's date
		var date string
		if len(args) > 5 {
			date = args[5]
		} else {
			date = time.Now().Format("2006-01-02")
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer db.Close()

		totalAmount := calculateTotalAmount(quantity, unitPrice)
		taxes := calculateTaxesFromTotalAmount(totalAmount)

		// Execute the insert query
		query := `INSERT INTO tblTransaction (id, stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(query, id, stockNo, date, quantity, tranType, unitPrice, totalAmount, taxes)
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
	},
}

func calculateTotalAmount(quantity int, unitPrice float64) float64 {
	return float64(quantity) * unitPrice
}

func calculateTaxesFromQuantityAndPrice(quantity int, unitPrice float64) int {
	var totalAmount float64 = calculateTotalAmount(quantity, unitPrice)
	return calculateTaxesFromTotalAmount(totalAmount)
}

func calculateTaxesFromTotalAmount(totalAmount float64) int {
	var taxRate float64 = 0.003
	return int(totalAmount * taxRate)
}
