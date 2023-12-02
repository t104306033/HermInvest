package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stock

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Stock management",
	Long:  `Manage stock via HermInestCli`,
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add stock",
	Long:  `Add stock to the inventory`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stock",
	Long:  `Delete stock from the inventory by providing the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update stock",
	Long:  `Update stock information, including unit price, in the inventory using the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query stock",
	Long:  `Query stock information from the inventory based on transaction ID, stock number, type, or date`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not implemented!")
	},
}

// version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The version of hermInvestCli is v0.0.1")
	},
}

// root

var rootCmd = &cobra.Command{
	Use:   "hermInvestCli",
	Short: "Oparate stock inventoy table", // working?
	Long:  `Oparate stock inventoy table for long desc`,
	Run: func(cmd *cobra.Command, args []string) {
		// if input is incorrect, show error and guide what to do
		// else if input is empty, show help
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
	rootCmd.AddCommand(versionCmd)

	stockCmd.AddCommand(addCmd)
	stockCmd.AddCommand(deleteCmd)
	stockCmd.AddCommand(updateCmd)
	stockCmd.AddCommand(queryCmd)
}

func main() {
	rootCmd.Execute()
}
