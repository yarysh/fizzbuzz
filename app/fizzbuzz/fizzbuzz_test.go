package fizzbuzz

import "testing"

func TestCalculate(t *testing.T) {
	tc := map[int64]string{
		-15: "FizzBuzz",
		-5:  "Buzz",
		-3:  "Fizz",
		-1:  "-1",
		0:   "FizzBuzz",
		1:   "1",
		3:   "Fizz",
		5:   "Buzz",
		15:  "FizzBuzz",
		33:  "Fizz",
		50:  "Buzz",
		150: "FizzBuzz",
	}
	for n, want := range tc {
		if Calculate(n) != want {
			t.Errorf("Calculate(%d) = %s, want %s", n, Calculate(n), want)
		}
	}
}
