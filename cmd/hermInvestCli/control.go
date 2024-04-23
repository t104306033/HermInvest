package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/repository"
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

func capitalReductionTransactionGenerator() {
	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	repo.DeleteAllTransactionRecordSys()

	// 1. Query all transaction records from tblCapitalReduction
	crs, _ := repo.QueryCapitalReductionAll()

	// 2. Iterate over each capital reduction record
	for _, cr := range crs {
		fmt.Println("\n---\n", cr)
		// Query transaction records by stock number
		trs, _ := repo.QueryTransactionRecordByStockNo(cr.StockNo, cr.CapitalReductionDate)

		remainingTrs := make([]*model.TransactionRecord, 0)
		for _, tr := range trs {
			if tr.TranType > 0 {
				remainingTrs = append(remainingTrs, tr)
			} else {
				qty := tr.Quantity
				for qty > 0 && len(remainingTrs) > 0 {
					var remove *model.TransactionRecord
					remove, remainingTrs = remainingTrs[0], remainingTrs[1:]
					qty -= remove.Quantity
				}
			}
		}
		var totalQuantity int
		var totalAmount int
		for _, tr := range remainingTrs {
			fmt.Println(tr)
			totalQuantity += tr.TranType * tr.Quantity
			totalAmount += int(float64(tr.Quantity) * tr.UnitPrice)
		}
		fmt.Println(totalQuantity, totalAmount)
		var avgUnitPrice float64 = float64(totalAmount) / float64(totalQuantity)
		fmt.Printf("%.2f\n", avgUnitPrice)

		// 3. insert into tblTransactionRecordSys
		capitalReductionRecord := model.NewTransactionRecord(
			cr.CapitalReductionDate, "08:00:00", cr.StockNo, -1, totalQuantity, avgUnitPrice)
		repo.InsertTransactionRecordSys(capitalReductionRecord)

		newStockNo := cr.NewStockNo
		if newStockNo == "" {
			newStockNo = cr.StockNo
		}

		newQuantity := int(float64(totalQuantity) * (1 - cr.Ratio))
		newAvgUnitPrice := (avgUnitPrice - cr.Cash) / (1 - cr.Ratio)
		newCapitalReductionRecord := model.NewTransactionRecord(
			cr.DistributionDate, "08:00:10", newStockNo, 1, newQuantity, newAvgUnitPrice)
		repo.InsertTransactionRecordSys(newCapitalReductionRecord)

		// deal with cash from tblCapitalReduction
	}

}

func transactionReGenerator() {
	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	serv := service.NewService(repo)

	// capitalReductionTransactionGenerator()

	repo.DeleteAlltblTransaction()
	repo.DeleteAlltblTransactionHistory()

	// repo.QueryTransactionRecordUnion()
	trs, _ := repo.QueryTransactionRecordUnion()

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
	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	newTransactions := make([]*model.Transaction, 0)
	newTransaction1 := model.NewTransactionFromInput(
		"2023-01-01", "09:00:00", "0050", 1, 2000, 22.6)
	newTransaction2 := model.NewTransactionFromInput(
		"2023-01-01", "09:00:00", "0050", 1, 3000, 22.7)
	newTransactions = append(newTransactions, newTransaction1)
	newTransactions = append(newTransactions, newTransaction2)

	ids, _ := repo.CreateTransactions(newTransactions)
	fmt.Print(ids)

}

func controlRun(cmd *cobra.Command, args []string) {
	capitalReductionTransactionGenerator()
}
