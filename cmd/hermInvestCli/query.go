package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query stock",
	Long:  `Query stock information from the inventory based on transaction ID, stock number, type, or date`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}
