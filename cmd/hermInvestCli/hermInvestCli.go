package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// stock

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Stock management",
	Long:  `Manage stock via HermInestCli`,
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stock",
	Long:  `Delete stock from the inventory by providing the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update stock",
	Long:  `Update stock information, including unit price, in the inventory using the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query stock",
	Long:  `Query stock information from the inventory based on transaction ID, stock number, type, or date`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

// version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The version of hermInvestCli is v0.0.5")
	},
}

// root

var rootCmd = &cobra.Command{
	Use:   "hermInvestCli",
	Short: "Oparate stock inventoy table", // working?
	Long:  `Oparate stock inventoy table for long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
	rootCmd.AddCommand(versionCmd)

	stockCmd.AddCommand(addCmd)
	stockCmd.AddCommand(deleteCmd)
	stockCmd.AddCommand(updateCmd)
	stockCmd.AddCommand(queryCmd)
}

func main() {
	rootCmd.Execute()
}
