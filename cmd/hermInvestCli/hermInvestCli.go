package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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

// version

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The version of hermInvestCli is v0.0.6")
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
