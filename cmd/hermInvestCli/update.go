package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. query transacation by id
// 3. recalc total amount and taxes
// 4. build sql syntax
// 5. insert into sql
// 6. print out result

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update stock",
	Long:  `Update stock information, including unit price, in the inventory using the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Please provide transaction ID and unit price")
			return
		}

		transactionID := args[0]
		unitPrice := args[1]

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
