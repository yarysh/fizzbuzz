package oracle

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	cacheMu sync.Mutex
	cache   map[int64]string
}

// FizzBuzz returns prediction of the fizzbuzz value for a given n
func (o *Oracle) FizzBuzz(n int64) (string, error) {
	o.cacheMu.Lock()
	if result, ok := o.cache[n]; ok {
		o.cacheMu.Unlock()
		log.Println("Get result from cache")
		return result, nil
	}
	o.cacheMu.Unlock()

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
	respResult := strings.Trim(string(body), "\"'\n")
	result := ""
	if respResult == fizzbuzz.Fizz || respResult == fizzbuzz.Buzz || respResult == fizzbuzz.FizzBuzz {
		result = respResult
	} else if _, err := strconv.ParseInt(respResult, 10, 64); err == nil {
		result = respResult
	}

	if result != "" {
		o.cacheMu.Lock()
		o.cache[n] = result
		o.cacheMu.Unlock()
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
		cache: map[int64]string{},
	}
}
