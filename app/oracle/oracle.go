package oracle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yarysh/fizzbuzz/app/fizzbuzz"
)

// Options for Oracle
type Options struct {
	BaseUrl string        // Base API url
	Timeout time.Duration // Timeout for API calls
}

// Oracle - prediction API
type Oracle struct {
	baseUrl string
	client  *http.Client
}

// FizzBuzz returns prediction of the fizzbuzz value for a given n
func (o *Oracle) FizzBuzz(n int64) (string, error) {
	path := "predict"
	resp, err := o.client.Post(
		o.baseUrl+path,
		"application/x-www-form-urlencoded",
		strings.NewReader(strconv.FormatInt(n, 10)),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Remove `"`, `'` and whitespaces
	result := strings.Trim(string(body), "\"'\n")
	if result == fizzbuzz.Fizz || result == fizzbuzz.Buzz || result == fizzbuzz.FizzBuzz {
		return result, nil
	} else if _, err := strconv.ParseInt(result, 10, 64); err == nil {
		return result, nil
	}
	return "", fmt.Errorf("unexpected result: %s", result)
}

// NewOracle returns a new Oracle
func NewOracle(opts Options) *Oracle {
	baseUrl := opts.BaseUrl
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	return &Oracle{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}
