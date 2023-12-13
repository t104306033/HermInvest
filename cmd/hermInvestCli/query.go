package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   `query {--all | --id <ID> | [--stockNo <StockNumber> --type <Type> --date <Date>]}`,
	Short: "Query stock (Transaction ID, Stock No., Type, or Date)",
	Example: "" +
		"  - Query by Transaction ID:\n" +
		"    hermInvestCli stock query --id 11\n\n" +

		"  - Query all records:\n" +
		"    hermInvestCli stock query --all\n\n" +

		"  - Query by stock number:\n" +
		"    hermInvestCli stock query --stockNo 0050\n\n" +

		"  - Query by stock number, type, and date:\n" +
		"    hermInvestCli stock query --stockNo 0050 --type 1 --date 2023-12-01",
	Long: "Query stock by transaction ID, stock number, type, or date.",
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
	queryCmd.Flags().Bool("all", false, "Query all records")
	queryCmd.Flags().Int("id", 0, "Query by ID")
	queryCmd.Flags().String("stockNo", "", "Stock number")
	queryCmd.Flags().Int("type", 0, "Type")
	queryCmd.Flags().String("date", "", "Date")
}
