package fizzbuzz

import "strconv"

const (
	Fizz     = "Fizz"
	Buzz     = "Buzz"
	FizzBuzz = "FizzBuzz"
)

// Calculate returns the fizzbuzz value for a given n
// see https://en.wikipedia.org/wiki/Fizz_buzz
func Calculate(n int64) string {
	switch {
	case n%3 == 0 && n%5 == 0:
		return FizzBuzz
	case n%3 == 0:
		return Fizz
	case n%5 == 0:
		return Buzz
	default:
		return strconv.FormatInt(n, 10)
	}
}
