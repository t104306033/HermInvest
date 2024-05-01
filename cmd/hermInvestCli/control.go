package main

import (
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

func controlRun(cmd *cobra.Command, args []string) {
	// capitalReductionTransactionGenerator()
	serv := service.InitializeService()

	// serv.RebuildCapitalReduction()
	err := serv.RebuildTransactionRecordSys()
	// err := serv.RebuildDividend()
	if err != nil {
		fmt.Println("error:", err)
	}
}
