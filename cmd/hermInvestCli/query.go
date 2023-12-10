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
		id, _ := cmd.Flags().GetInt("id")
		stockNo, _ := cmd.Flags().GetString("stockNo")
		tranType, _ := cmd.Flags().GetInt("type")
		date, _ := cmd.Flags().GetString("date")

		if all {
			query = buildQueryAll()
		} else if id != 0 {
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
	return fmt.Sprintln("SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction")
}
func buildQueryByID(id int) string {
	return fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE id = '%d'", id)
}

func buildQueryByDetails(stockNo string, tranType int, date string) string {
	var conditions []string

	if stockNo != "" {
		conditions = append(conditions, fmt.Sprintf("stockNo = '%s'", stockNo))
	}
	if tranType != 0 {
		conditions = append(conditions, fmt.Sprintf("tranType = '%d'", tranType))
	}
	if date != "" {
		conditions = append(conditions, fmt.Sprintf("date = '%s'", date))
	}

	var query string
	if len(conditions) > 0 {
		query = fmt.Sprintf("SELECT id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes FROM tblTransaction WHERE %s", strings.Join(conditions, " AND "))
	} else {
		query = "SELECT id, stockNo, tranType, quantity, unitPrice FROM tblTransaction"
	}

	return query
}

func displayResults(rows *sql.Rows) {
	fmt.Print("ID,\tStock No,\tType,\tQty(shares),\tUnit Price,\tTotal Amount,\ttaxes\n")
	for rows.Next() {
		var id int
		var stockNo string
		// var stockName string
		var tranType, quantity int
		var unitPrice float64
		// var date string
		var totalAmount, taxes int

		err := rows.Scan(&id, &stockNo, &tranType, &quantity, &unitPrice, &totalAmount, &taxes)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		fmt.Printf("%d,\t%s,\t\t%d,\t%d,\t\t%.2f,\t\t%d,\t\t%d\n", id, stockNo, tranType, quantity, unitPrice, totalAmount, taxes)
	}
}

func init() {
	queryCmd.Flags().Bool("all", false, "query all. indepent")
	queryCmd.Flags().Int("id", 0, "query by id. indepent")
	queryCmd.Flags().String("stockNo", "", "Stock number")
	queryCmd.Flags().Int("type", 0, "Type")
	queryCmd.Flags().String("date", "", "Date")
}
