package fizzbuzz

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OracleOptions struct {
	BaseUrl string        // Base API url
	Timeout time.Duration // Timeout for API calls
}

// Oracle - FizzBuzz prediction API
type Oracle struct {
	baseUrl string
	client  *http.Client
}

// Predict returns prediction of the fizzbuzz value for a given n
func (o *Oracle) Predict(n int64) (string, error) {
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

	result := strings.TrimSpace(string(body))
	if result == Fizz || result == Buzz || result == FizzBuzz {
		return result, nil
	} else if _, err := strconv.ParseInt(result, 10, 64); err == nil {
		return result, nil
	}
	return "", fmt.Errorf("unexpected result: %s", result)
}

// NewOracle returns a new Oracle
func NewOracle(opts OracleOptions) *Oracle {
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
