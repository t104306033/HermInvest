package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update stock",
	Long:  `Update stock information, including unit price, in the inventory using the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}
