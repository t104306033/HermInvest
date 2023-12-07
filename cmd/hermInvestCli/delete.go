package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stock",
	Long:  `Delete stock from the inventory by providing the transaction ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the transaction ID to delete.")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid ID provided. Please provide a valid ID.")
			return
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		defer db.Close()

		confirm := confirmDeletion()
		if confirm {
			err = deleteTransaction(db, id)
			if err != nil {
				fmt.Println("Error deleting transaction:", err)
				return
			}
			fmt.Println("Transaction deleted successfully!")
		} else {
			fmt.Println("Deletion cancelled.")
		}
	},
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

func deleteTransaction(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tblTransaction WHERE id = ?", id)
	return err
}
