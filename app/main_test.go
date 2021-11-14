package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewFizzBuzzService(t *testing.T) {
	server := httptest.NewServer(NewFizzBuzzService(""))
	defer server.Close()

	// GET request should return status "405 Method Not Allowed"
	{
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		client := &http.Client{Timeout: time.Second}
		resp, _ := client.Do(req)
		want := "405 Method Not Allowed"
		if resp.Status != want {
			t.Errorf("Expected status %s, got %s", want, resp.Status)
		}
	}

	// POST request should return status "200 OK" and body "FizzBuzz"
	{
		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer([]byte("15")))
		client := &http.Client{Timeout: time.Second}
		resp, _ := client.Do(req)
		want := "200 OK"
		if resp.Status != want {
			t.Errorf("Expected status %s, got %s", want, resp.Status)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		want = `"FizzBuzz"`
		if string(body) != want {
			t.Errorf("Expected body %s, got %s", want, string(body))
		}
	}
}
