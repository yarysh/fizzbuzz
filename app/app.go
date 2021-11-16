package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/yarysh/fizzbuzz/app/fizzbuzz"
)

type PredictAPI interface {
	FizzBuzz(n int64) (string, error)
}

// App - FizzBuzz application
// Values in the range [LocalCalcRange[0], LocalCalcRange[1]] (inclusive)
// will be calculated locally, others will be using Oracle Prediction API.
type App struct {
	Oracle         PredictAPI
	LocalCalcRange [2]int64
}

// formatResult prepares result for output
func formatResult(result string) []byte {
	return []byte(fmt.Sprintf("%q", result))
}

// FizzBuzzHandler handles requests to calculate fizzbuzz
func (a *App) FizzBuzzHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Couldn't read request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body."))
		return
	}

	data := strings.TrimSpace(string(rawData))
	n, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Printf("Couldn't convert %q to int64: %v\n", data, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request body malformed."))
		return
	}

	// For n between LocalCalcRange[0] and LocalCalcRange[1] (inclusive), use local FizzBuzz calculation
	if n >= a.LocalCalcRange[0] && n <= a.LocalCalcRange[1] {
		w.Write(formatResult(fizzbuzz.Calculate(n)))
		return
	}
	// In other cases, use Oracle prediction API
	result, err := a.Oracle.FizzBuzz(n)
	log.Printf("Oracle API was requested with n = %d\n", n)
	if err != nil {
		log.Printf("Oracle API coulnd't predict result for value %d: %v\n", n, err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Prediction API unavailable or returned an error."))
		return
	}
	w.Write(formatResult(result))
}
