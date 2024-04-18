package main

import (
	"HermInvest/pkg/repository"
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
	defer db.Close()

	// init transactionRepository
	repo := repository.NewTransactionRepository(db)

	// 1. select * from tblCapitalReduction
	crs, err := repo.QueryCapitalReductionAll()

	// 2. select stockNo quantity tblTransactionRecord group by and where stockNo
	repo.QueryTransactionRecordByStockNo(crs[1].StockNo)
	// 3. insert into tblTransactionRecordSys
}

func controlRun(cmd *cobra.Command, args []string) {
	capitalReductionTransactionGenerator()
}
