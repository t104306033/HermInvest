package main

import (
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
	// capitalReductionTransactionGenerator()
	// serv := service.InitializeService()

	http.HandleFunc("/", homePage)

	open("http://127.0.0.1:9453")

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
