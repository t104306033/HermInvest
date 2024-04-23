package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/service"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. calc total amount and taxes
// 3. build sql syntax
// 4. insert into sql
// 5. print out result

var addCmd = &cobra.Command{
	Use:   "add date time stockNo type quantity unitPrice",
	Short: "Add date time stock (Stock No., Type, Quantity, Unit Price)",
	Example: "" +
		"  - Purchase on a specific date:\n" +
		"    hermInvestCli stock add 2023-12-01 09:00:00 0050 1 1500 23.5\n\n" +

		"  - Sale on a specific date:\n" +
		"    hermInvestCli stock add -- 2023-12-01 09:00:10 0050 -1 1500 23.5",
	Long: `Add stock by transaction date time stockNo type quantity unitPrice`,
	Args: cobra.RangeArgs(6, 6),
	Run:  addRun,
}

func init() {
	stockCmd.AddCommand(addCmd)
}

func addRun(cmd *cobra.Command, args []string) {
	tranDate, tranTime, stockNo, tranType, quantity, unitPrice, err := ParseTransactionForAddCmd(args)
	if err != nil {
		fmt.Println("Error parsing transaction data:", err)
		return
	}

	serv := service.InitializeService()

	// add stock in inventory
	// 1. new transaction from input
	// 2. find the first purchase from the inventory
	// 3. check the transaction type of new transaction and first purchase

	// TODO: service.addTransaction() AddTransactionAndUpdateInventory
	newTransaction := model.NewTransactionFromInput(tranDate, tranTime, stockNo, tranType, quantity, unitPrice)

	t, err := serv.AddTransaction(newTransaction)
	if err != nil {
		fmt.Println("Error adding transaction: ", err)
	} else if t != nil {
		var ts []*model.Transaction
		ts = append(ts, t)
		displayResults(ts)
	}

}

func ParseTransactionForAddCmd(args []string) (string, string, string, int, int, float64, error) {

	var parsedTime time.Time
	var err error
	parsedTime, err = time.Parse(time.DateOnly, args[0])
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("error parsing date: %s", err)
	}
	tranDate := parsedTime.Format(time.DateOnly)

	parsedTime, err = time.Parse(time.TimeOnly, args[1])
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("error parsing time: %s", err)
	}
	tranTate := parsedTime.Format(time.TimeOnly)

	stockNo := args[2] // regex a-z A-Z 0-9

	tranType, err := strconv.Atoi(args[3])
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("error parsing integer: %s", err)
	}

	quantity, err := strconv.Atoi(args[4])
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("error parsing integer: %s", err)
	}

	unitPrice, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return "", "", "", 0, 0, 0, fmt.Errorf("error parsing float: %s", err)
	}

	return tranDate, tranTate, stockNo, tranType, quantity, unitPrice, nil
}
