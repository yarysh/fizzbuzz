package main

import (
	"bytes"
	"fmt"
	"github.com/yarysh/fizzbuzz/app/oracle"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_FizzBuzzHandler(t *testing.T) {
	orcl, mux, teardown := oracle.TestOracle(t)
	defer teardown()
	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		switch string(body) {
		case "-1":
			fmt.Fprint(w, `"Buzz"`)
		case "101":
			fmt.Fprint(w, `"FizzBuzz"`)
		case "150":
			fmt.Fprint(w, `"Fizz"`)
		default:
			t.Errorf("Unexpected request: %s", string(body))
		}
	})
	tc := map[string]string{
		"-1":  `"Buzz"`,     // Oracle prediction (wrong result to detect that calculation was done by Oracle)
		"1":   `"1"`,        // Local calculation
		"3":   `"Fizz"`,     // Local calculation
		"5":   `"Buzz"`,     // Local calculation
		"15":  `"FizzBuzz"`, // Local calculation
		"100": `"Buzz"`,     // Local calculation
		"101": `"FizzBuzz"`, // Oracle calculation (wrong result to detect that calculation was done by Oracle)
		"150": `"Fizz"`,     // Oracle prediction (wrong result to detect that calculation was done by Oracle)
	}
	app := App{
		Oracle:         orcl,
		LocalCalcRange: [2]int64{1, 100},
	}
	for n, want := range tc {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(n)))
		app.FizzBuzzHandler(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		got, _ := ioutil.ReadAll(resp.Body)
		if string(got) != want {
			t.Errorf("FizzBuzzHandler(%s) returns %s, want %s", n, string(got), want)
		}
	}
}