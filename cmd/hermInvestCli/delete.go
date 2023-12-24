package main

import (
	"HermInvest/pkg/repository"
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete id",
	Short: "Delete stock by transaction ID",
	Example: "" +
		"  - Delete by ID:\n" +
		"    hermInvestCli stock delete 11",
	Long: `Delete stock from the inventory by providing the stock transaction ID`,
	Args: cobra.ExactArgs(1),
	Run:  deleteRun,
}

func init() {
	stockCmd.AddCommand(deleteCmd)
}

func deleteRun(cmd *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid ID provided. Please provide a valid ID.")
		return
	}

	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}
	defer db.Close()

	// init transactionRepository
	repo := &repository.TransactionRepository{DB: db}

	confirm := confirmDeletion()
	if confirm {
		err = repo.DeleteTransaction(id)
		if err != nil {
			fmt.Println("Error deleting transaction:", err)
			return
		}
		fmt.Println("Transaction deleted successfully!")
	} else {
		fmt.Println("Deletion cancelled.")
	}
}

func confirmDeletion() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you want to delete this transaction? (yes/no): ")
	text, _ := reader.ReadString('\n')
	text = trimNewline(text)

	return text == "yes"
}

func trimNewline(s string) string {
	return s[:len(s)-1]
}
