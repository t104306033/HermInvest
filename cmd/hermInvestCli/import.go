package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/repository"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
		"    hermInvestCli stock import stock.csv\n\n" +

		"  - Import stock from file and swap new order column as 0,1,3,2,4:\n" +
		"    hermInvestCli stock import stock.csv --swapColumn 0,1,3,2",
	Long: "" +
		"Import stock from csv file.\n" +
		"Please check your csv file has column stockNo type quantity unitPrice.",
	Args: cobra.ExactArgs(1),
	Run:  importRun,
}

func init() {
	stockCmd.AddCommand(importCmd)

	importCmd.Flags().Bool("skipHeader", false, "Ignore header")
	importCmd.Flags().String("swapColumn", "", "Swap column")
}

func importRun(cmd *cobra.Command, args []string) {
	filePath := args[0]
	skipHeader, _ := cmd.Flags().GetBool("skipHeader")
	indexes, _ := cmd.Flags().GetString("swapColumn")

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
	repo := &repository.TransactionRepository{DB: db}

	var transactions []*model.Transaction
	for _, row := range rows {
		if indexes != "" {
			row, err = swapColumn(row, indexes)
			if err != nil {
				fmt.Println("Error swaping column:", err)
				return
			}
		}

		stockNo, tranType, quantity, unitPrice, date, err := ParseTransactionForAddCmd(row)
		if err != nil {
			fmt.Println("Error parsing transaction data:", err)
			return
		}

		t := model.NewTransactionFromInput(stockNo, date, quantity, tranType, unitPrice)
		transactions = append(transactions, t)
	}

	// TODO: create Transactions, bulk insert? Finally, I choose begin a db transaction
	ids, err := repo.CreateTransactions(transactions)
	if err != nil {
		fmt.Println("Error creating transaction: ", err)
	}

	fmt.Println("inserted ids:", ids)

	// 	// Bool control show result or not
	// 	transactions, err := repo.queryTransactionByID(id)
	// 	if err != nil {
	// 		fmt.Println("Error querying database:", err)
	// 	}

	// Print out result
	// displayResults(transactions)
}

func swapColumn(row []string, indexes string) ([]string, error) {
	if indexes == "" {
		return row, errors.New("indexes cannot be empty")
	}

	indexArr := strings.Split(indexes, ",")
	if len(indexArr) > len(row) {
		return nil, errors.New("indexes length exceeds row length")
	}

	rowForSwapping := make([]string, len(row))
	copy(rowForSwapping, row)

	for i, idxStr := range indexArr {
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return nil, fmt.Errorf("parsing index '%s': %w", idxStr, err)
		}
		if idx < 0 || idx >= len(row) {
			return nil, fmt.Errorf("index '%d' is out of range", idx)
		}

		row[i] = rowForSwapping[idx]
	}

	return row, nil
}
