package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stock",
	Long:  `Delete stock from the inventory by providing the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}
