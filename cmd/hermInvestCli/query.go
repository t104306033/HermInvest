package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query stock",
	Long: `Query stock information from the inventory based on transaction ID, stock number, type, or date.

	Usage:
	  hermInvestCli stock query [stock number] [type] [date] [flags]
	  hermInvestCli stock query byID [transaction ID] [flags]`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 && len(args) == 0 {
			fmt.Println("No args and no flags provided.")
			cmd.Help()
			return
		}

		var query string

		all, _ := cmd.Flags().GetBool("all")
		id, _ := cmd.Flags().GetString("id")
		stockNo, _ := cmd.Flags().GetString("stockNo")
		tranType, _ := cmd.Flags().GetString("type")
		date, _ := cmd.Flags().GetString("date")

		if all {
			query = buildQueryAll()
		} else if id != "" {
			query = buildQueryByID(id)
		} else {
			query = buildQueryByDetails(stockNo, tranType, date)
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer db.Close()

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println("Error querying database:", err)
			return
		}
		defer rows.Close()

		displayResults(rows)
	},
}

func buildQueryAll() string {
	return fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice FROM tblTransaction")
}
func buildQueryByID(id string) string {
	return fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice FROM tblTransaction WHERE id = '%s'", id)
}

func buildQueryByDetails(stockNo, tranType, date string) string {
	var conditions []string

	if stockNo != "" {
		conditions = append(conditions, fmt.Sprintf("stockNo = '%s'", stockNo))
	}
	if tranType != "" {
		conditions = append(conditions, fmt.Sprintf("tranType = '%s'", tranType))
	}
	if date != "" {
		conditions = append(conditions, fmt.Sprintf("date = '%s'", date))
	}

	var query string
	if len(conditions) > 0 {
		query = fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))
	} else {
		query = "SELECT id, stockNo, tranType, quantity, unitPrice FROM tblTransaction"
	}

	return query
}

func displayResults(rows *sql.Rows) {
	for rows.Next() {
		var id int
		var stockNo string
		// var stockName string
		var tranType, quantity int
		var unitPrice float64
		// var date string
		// var totalAmount, taxes int

		if err := rows.Scan(&id, &stockNo, &tranType, &quantity, &unitPrice); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		fmt.Printf("ID: %d, Stock No: %s, Type: %d, Quantity: %d, Unit Price: %.2f\n", id, stockNo, tranType, quantity, unitPrice)

	}
}

func init() {
	queryCmd.Flags().Bool("all", false, "query all. indepent")
	queryCmd.Flags().String("id", "", "query by id. indepent")
	queryCmd.Flags().String("stockNo", "", "Stock number")
	queryCmd.Flags().String("type", "", "Type")
	queryCmd.Flags().String("date", "", "Date")
}
