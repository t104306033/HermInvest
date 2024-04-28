package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/service"
	"fmt"

	"github.com/spf13/cobra"
)

var controlCmd = &cobra.Command{
	Use: "control",
	Run: controlRun,
}

func init() {
	stockCmd.AddCommand(controlCmd)
}

func transactionReGenerator() {
	serv := service.InitializeService()

	// capitalReductionTransactionGenerator()

	serv.DeleteAlltblTransaction()
	serv.DeleteAlltblTransactionHistory()

	// repo.QueryTransactionRecordUnion()
	trs, _ := serv.QueryTransactionRecordUnion()

	var transactions []*model.Transaction
	for _, tr := range trs {
		newTransaction := model.NewTransactionFromInput(
			tr.Date, tr.Time, tr.StockNo, tr.TranType, tr.Quantity, tr.UnitPrice)
		t, err := serv.AddTransaction(newTransaction)
		if err != nil {
			fmt.Println("Error adding transaction: ", err)
		} else if t != nil {
			transactions = append(transactions, t)
		}
	}
	displayResults(transactions)
}

func test() {
	serv := service.InitializeService()

	newTransactions := make([]*model.Transaction, 0)
	newTransaction1 := model.NewTransactionFromInput(
		"2023-01-01", "09:00:00", "0050", 1, 2000, 22.6)
	newTransaction2 := model.NewTransactionFromInput(
		"2023-01-01", "09:00:00", "0050", 1, 3000, 22.7)
	newTransactions = append(newTransactions, newTransaction1)
	newTransactions = append(newTransactions, newTransaction2)

	ids, _ := serv.CreateTransactions(newTransactions)
	fmt.Print(ids)

}

func controlRun(cmd *cobra.Command, args []string) {
	// capitalReductionTransactionGenerator()
	serv := service.InitializeService()

	serv.RebuildCapitalReduction()
}
