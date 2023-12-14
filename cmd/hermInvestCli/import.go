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
		skipHeader, _ := cmd.Flags().GetBool("skipHeader")
		filePath := args[0]

		// check file path exist
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Error filePath not exist: ", err)
			} else {
				fmt.Println("Error getting file info: ", err)
			}
			return
		}

		// check it is not a dir
		if fileInfo.IsDir() {
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

		for _, row := range rows {
			fmt.Println(row)
		}
	},
}

func init() {
	importCmd.Flags().Bool("skipHeader", false, "Ignore header in CSV")
}
