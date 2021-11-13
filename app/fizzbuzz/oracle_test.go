package fizzbuzz

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestOracle_Predict(t *testing.T) {
	oracle, mux, teardown := TestOracle(t)
	defer teardown()

	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		switch string(body) {
		case "1":
			fmt.Fprint(w, "1")
		case "3":
			fmt.Fprint(w, "Fizz")
		case "5":
			fmt.Fprint(w, "Buzz")
		case "15":
			fmt.Fprint(w, "FizzBuzz")
		}
	})

	tc := map[int64]string{
		1:  "1",
		3:  "Fizz",
		5:  "Buzz",
		15: "FizzBuzz",
	}
	for n, want := range tc {
		if got, _ := oracle.Predict(n); got != want {
			t.Errorf("Predict(%d) = %s, want %s", n, got, want)
		}
	}
}

func TestOracle_Predict_bad_response(t *testing.T) {
	oracle, mux, teardown := TestOracle(t)
	defer teardown()

	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
	})
	_, err := oracle.Predict(1)
	want := "unexpected response: 500 Internal Server Error"
	if err.Error() != want {
		t.Errorf("expected error %q, got %q", want, err.Error())
	}
}

func TestOracle_Predict_timeout(t *testing.T) {
	oracle, mux, teardown := TestOracle(t)
	defer teardown()

	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		fmt.Fprint(w, "Fizz")
	})

	_, err := oracle.Predict(1)
	if !strings.Contains(err.Error(), "Client.Timeout exceeded") {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestOracle_Predict_unknown_result(t *testing.T) {
	oracle, mux, teardown := TestOracle(t)
	defer teardown()

	mux.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Fizz1")
	})
	_, err := oracle.Predict(1)
	want := "unexpected result: Fizz1"
	if err.Error() != want {
		t.Errorf("expected error: %s, got: %v", want, err)
	}
}
