package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// 1. check input
// 2. calc total amount and taxes
// 3. build sql syntax
// 4. insert into sql
// 5. success

// TODO: ignore header (first line)
// TODO: select stockNo, quantity ... from csv
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
		filePath := args[0]

		// check input is file
		// check file path exist and not a dir
		fileInfo, err := os.Stat(filePath) // os.Open also seems to have Stat()
		if os.IsNotExist(err) {
			fmt.Println("Error filePath not exist: ", err)
			return
		} else if fileInfo.IsDir() {
			fmt.Println("Error filePath is not a file.")
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening the file: ", err)
			return
		}
		defer file.Close()

		fileReader := csv.NewReader(file)

		rows, err := fileReader.ReadAll()
		if err != nil {
			fmt.Println("Error reading rows: ")
			return
		}

		for _, row := range rows {
			fmt.Println(row[0])
		}
	},
}
