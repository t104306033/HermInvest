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
	db, err := repository.GetDBConnectionGorm()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}

	// init transactionRepository
	repo := repository.NewTransactionRepositoryGorm(db)

	// 1. select * from tblCapitalReduction
	repo.QueryCapitalReductionAll()

	// 2. select stockNo quantity tblTransactionRecord group by and where stockNo
	repo.QueryTransactionRecordByStockNo("2409")
	// 3. insert into tblTransactionRecordSys
}

func controlRun(cmd *cobra.Command, args []string) {
	capitalReductionTransactionGenerator()
}
