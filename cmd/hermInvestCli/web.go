package main

import (
	"HermInvest/pkg/model"
	"HermInvest/pkg/service"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use: "web",
	Run: webRun,
}

func init() {
	stockCmd.AddCommand(webCmd)
}

func webRun(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/transaction", transactionPage)

	open("http://127.0.0.1:9453/transaction")

	err := http.ListenAndServe(":9453", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homePage(w http.ResponseWriter, _ *http.Request) {
	var pageHTML []byte
	pageHTML, err := os.ReadFile("html/home.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	w.Write([]byte(pageHTML))
}

func transactionPage(w http.ResponseWriter, _ *http.Request) {
	serv := service.InitializeService()

	transactions, err := serv.QueryTransactionAll()
	if err != nil {
		fmt.Println("Error querying database:", err)
	}
	// else {
	// 	displayResults(transactions) // for loop print transactions
	// }

	tableHTML := displayResultsHMTL(transactions)

	var pageHTML []byte
	pageHTML, err = os.ReadFile("html/transaction.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	// Split the HTML byte slice into two parts
	var firstPart, secondPart []byte
	splitIndex := bytes.Index(pageHTML, []byte("<TransactionTable>"))
	if splitIndex != -1 {
		firstPart = pageHTML[:splitIndex]
		secondPart = pageHTML[splitIndex+len("<TransactionTable>"):]
	} else {
		// Handle the case where the "here" marker is not found
		log.Fatal("Marker '<TransactionTable>' not found in HTML template")
	}

	// Insert the transaction table HTML between the two parts
	combinedHTML := append(firstPart, tableHTML...)
	combinedHTML = append(combinedHTML, secondPart...)

	w.Write([]byte(combinedHTML))
}

func displayResultsHMTL(transactions []*model.Transaction) string {
	// Create a table HTML string
	tableHTML := "<table>"
	tableHTML += "<tr><th>ID</th><th>Stock No</th><th>Type</th><th>Qty(shares)</th><th>Unit Price</th><th>Total Amount</th><th>Taxes</th></tr>"

	// Iterate through transactions and add rows to the table
	for _, transaction := range transactions {
		tableHTML += "<tr>"
		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.ID)
		tableHTML += fmt.Sprintf("<td>%s</td>", transaction.StockNo)
		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.TranType)
		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.Quantity)
		tableHTML += fmt.Sprintf("<td>%.2f</td>", transaction.UnitPrice)
		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.TotalAmount)
		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.Taxes)
		tableHTML += "</tr>"
	}

	tableHTML += "</table>"

	// Print the table HTML to the console for debugging (remove this in production)
	// fmt.Println(tableHTML)

	return tableHTML
}

func open(url string) error { // open url from browser
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
