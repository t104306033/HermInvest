package main

import (
	"HermInvest/pkg/service"

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
	serv.RebuildTransaction()
}
