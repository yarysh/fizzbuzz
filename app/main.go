package main

import (
	"flag"
	"fmt"
	"github.com/yarysh/fizzbuzz/app/oracle"
	"net/http"
	"os"
)

func NewFizzBuzzService(oracleUrl string) *Server {
	app := App{
		Oracle: oracle.NewOracle(oracle.Options{
			BaseUrl: oracleUrl,
		}),
		LocalCalcRange: [2]int64{1, 100},
	}
	server := NewServer()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only POST requests are allowed
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		app.FizzBuzzHandler(w, r)
	})

	return server
}

func main() {
	var addr string
	var oracleUrl string
	flag.StringVar(&addr, "addr", ":8080", "server addr")
	flag.StringVar(&oracleUrl, "oracle_url", "https://fizzbuzz-oracle.b17g.services/", "oracle api url")
	flag.Parse()

	fmt.Printf("Starting server at %s...\n", addr)
	err := http.ListenAndServe(addr, NewFizzBuzzService(oracleUrl))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
