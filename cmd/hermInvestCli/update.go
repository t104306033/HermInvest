package main

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

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

		// TODO: extract DB connection to configuration
		// SQLite3 connection with foreign keys enabled
		dbPath := "./internal/app/database/dev-database.db?_foreign_keys=true"
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Println("Error connecting to the database:", err)
			return
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
