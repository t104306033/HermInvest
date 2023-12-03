package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add id stockNo type quantity unitPrice [date]",
	Short: "Add stock",
	Long:  `Add stock to the inventory`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		stockNo := args[1]
		tranType := args[2]
		quantity := args[3]
		unitPrice := args[4]

		// Parse date argument or default is today's date
		var date string
		if len(args) > 5 {
			date = args[5]
		} else {
			date = time.Now().Format("2006-01-02")
		}

		// sqlite3 connection with  foreign keys enabled
		var dbPath = "./internal/app/database/dev-database.db?_foreign_keys=true"

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer db.Close()

		// Execute the insert query
		query := `INSERT INTO tblTransaction (id, stockNo, date, quantity, tranType, unitPrice) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(query, id, stockNo, date, quantity, tranType, unitPrice)
		if err != nil {
			fmt.Println("Error: ", err)
		} else {
			fmt.Println("Pass: Stock added successfully!")
		}
	},
}
