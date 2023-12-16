package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. calc total amount and taxes
// 3. build sql syntax
// 4. insert into sql
// 5. success

// TODO: column format

var importCmd = &cobra.Command{
	Use:   "import file",
	Short: "Import stock from csv file",
	Example: "" +
		"  - Import stock from file:\n" +
		"    hermInvestCli stock import stock.csv",
	Long: "" +
		"Import stock from csv file.\n" +
		"Please check your csv file has column stockNo type quantity unitPrice.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		skipHeader, _ := cmd.Flags().GetBool("skipHeader")
		filePath := args[0]

		// need testcase check file path exist
		// need testcase check file permission
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening the file: ", err)
			return
		}
		defer file.Close()

		// need testcase check it is not a dir
		fileInfo, _ := file.Stat()
		if fileInfo.IsDir() {
			fmt.Println("Error filePath is not a file.")
			return
		}

		fileReader := csv.NewReader(file)

		if skipHeader {
			_, err := fileReader.Read()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Error file empty.")
				} else {
					fmt.Println("Error reading header: ", err)
				}
				return
			}
		}

		rows, err := fileReader.ReadAll()
		if err != nil {
			fmt.Println("Error reading rows: ", err)
			return
		}

		db, err := GetDBConnection()
		if err != nil {
			fmt.Println("Error geting DB connection: ", err)
		}
		defer db.Close()

		// init transactionRepository
		repo := &transactionRepository{db: db}

		var transactions []*Transaction
		for _, row := range rows {
			// TODO: select stockNo, quantity ... from csv (swap row)
			stockNo, tranType, quantity, unitPrice, date, err := ParseTransactionForAddCmd(row)
			if err != nil {
				fmt.Println("Error parsing transaction data:", err)
				return
			}

			t := newTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)
			transactions = append(transactions, t)
		}

		// TODO: create Transactions, bulk insert? Finally, I choose begin a db transaction
		err = repo.createTransactions(transactions)
		if err != nil {
			fmt.Println("Error creating transaction: ", err)
		}

		// 	// Bool control show result or not
		// 	transactions, err := repo.queryTransactionByID(id)
		// 	if err != nil {
		// 		fmt.Println("Error querying database:", err)
		// 	}

		// Print out result
		// displayResults(transactions)
	},
}

func init() {
	importCmd.Flags().Bool("skipHeader", false, "Ignore header in CSV")
}
