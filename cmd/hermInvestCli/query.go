package main

import (
	"HermInvest/pkg/model"
	"fmt"

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
	Args: cobra.NoArgs,
	RunE: queryRun,
}

func init() {
	stockCmd.AddCommand(queryCmd)

	queryCmd.Flags().Bool("all", false, "Query all records")
	queryCmd.Flags().Int("id", 0, "Query by ID")
	queryCmd.Flags().String("stockNo", "", "Stock number")
	queryCmd.Flags().Int("type", 0, "Type")
	queryCmd.Flags().String("date", "", "Date")
}

func queryRun(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 {
		return fmt.Errorf("no flags provided")
	}

	all, _ := cmd.Flags().GetBool("all")
	id, _ := cmd.Flags().GetInt("id")
	stockNo, _ := cmd.Flags().GetString("stockNo")
	tranType, _ := cmd.Flags().GetInt("type")
	date, _ := cmd.Flags().GetString("date")

	db, err := GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}
	defer db.Close()

	// init transactionRepository
	repo := &transactionRepository{db: db}

	var transactions []*model.Transaction
	var transactionsErr error
	if all {
		transactions, transactionsErr = repo.queryTransactionAll()
	} else if id != 0 {
		transactions, transactionsErr = repo.queryTransactionByID(id)
	} else {
		transactions, transactionsErr = repo.queryTransactionByDetails(stockNo, tranType, date)
	}
	if transactionsErr != nil {
		fmt.Println("Error querying database:", transactionsErr)
	}

	displayResults(transactions)

	return nil
}

func displayResults(transactions []*model.Transaction) {
	fmt.Print("ID,\tStock No,\tType,\tQty(shares),\tUnit Price,\tTotal Amount,\ttaxes\n")
	for _, t := range transactions {
		fmt.Printf("%d,\t%s,\t\t%d,\t%d,\t\t%.2f,\t\t%d,\t\t%d\n", t.ID, t.StockNo, t.TranType, t.Quantity, t.UnitPrice, t.TotalAmount, t.Taxes)
	}
}
