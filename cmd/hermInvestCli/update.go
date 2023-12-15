package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. query transacation by id
// 3. recalc total amount and taxes
// 4. build sql syntax
// 5. insert into sql
// 6. print out result

var updateCmd = &cobra.Command{
	Use:   "update id unitPrice",
	Short: "Update unit price by transaction ID",
	Example: "" +
		"  - Update unit Price by ID:\n" +
		"    hermInvestCli stock update 11 20.3",
	Long: `Update the unit price of stock in the inventory using the transaction ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Please provide transaction ID and unit price")
			cmd.Help()
			return
		}

		transactionID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error parsing integer: ", err)
		}
		unitPrice, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println("Error parsing float: ", err)
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer db.Close()

		// TODO: Recalculate totalAmount, taxes
		// TODO: check update work. ex: update a fake transaction ID to db
		// Execute update query
		query := "UPDATE tblTransaction SET unitPrice = ? WHERE id = ?"
		_, err = db.Exec(query, unitPrice, transactionID)
		if err != nil {
			fmt.Println("Error updating stock information:", err)
			return
		}

		fmt.Printf("Successfully updated transaction ID %s with new unit price %s\n", transactionID, unitPrice)
	},
}
