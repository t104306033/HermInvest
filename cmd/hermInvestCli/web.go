package main

import (
	"HermInvest/pkg/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
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
	router := gin.Default()

	router.GET("/", homePage)
	router.GET("/transaction", transactionPage)
	router.GET("/api/transaction", apiGetTransactions)

	open("http://127.0.0.1:9453/transaction")

	err := router.Run(":9453")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homePage(c *gin.Context) {
	var pageHTML []byte
	pageHTML, err := os.ReadFile("html/home.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	c.Data(http.StatusOK, "text/html", pageHTML)
}

func apiGetTransactions(c *gin.Context) {
	// transactions := []*model.Transaction{
	// 	{ID: 1, StockNo: "ABC", TranType: 1, Quantity: 100, UnitPrice: 10.50, TotalAmount: 1050, Taxes: 50},
	// 	{ID: 2, StockNo: "XYZ", TranType: 2, Quantity: 50, UnitPrice: 20.25, TotalAmount: 1012, Taxes: 12},
	// }
	serv := service.InitializeService()

	transactions, err := serv.QueryTransactionAll()
	if err != nil {
		fmt.Println("Error querying database:", err)
	}

	c.JSON(http.StatusOK, transactions)
}

func transactionPage(c *gin.Context) {

	var pageHTML []byte
	pageHTML, err := os.ReadFile("html/transaction.html")
	if err != nil {
		log.Fatal("os.ReadFile: ", err)
	}

	c.Data(http.StatusOK, "text/html", pageHTML)
}

// func displayResultsHMTL(transactions []*model.Transaction) string {
// 	// Create a table HTML string
// 	tableHTML := "<table>"
// 	tableHTML += "<tr><th>ID</th><th>Stock No</th><th>Type</th><th>Qty(shares)</th><th>Unit Price</th><th>Total Amount</th><th>Taxes</th></tr>"

// 	// Iterate through transactions and add rows to the table
// 	for _, transaction := range transactions {
// 		tableHTML += "<tr>"
// 		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.ID)
// 		tableHTML += fmt.Sprintf("<td>%s</td>", transaction.StockNo)
// 		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.TranType)
// 		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.Quantity)
// 		tableHTML += fmt.Sprintf("<td>%.2f</td>", transaction.UnitPrice)
// 		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.TotalAmount)
// 		tableHTML += fmt.Sprintf("<td>%d</td>", transaction.Taxes)
// 		tableHTML += "</tr>"
// 	}

// 	tableHTML += "</table>"

// 	// Print the table HTML to the console for debugging (remove this in production)
// 	// fmt.Println(tableHTML)

// 	return tableHTML
// }

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
