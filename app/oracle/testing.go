package oracle

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestOracle - Oracle test server
func TestOracle(t *testing.T) (*Oracle, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	oracle := NewOracle(Options{
		BaseUrl: server.URL,
		Timeout: 100 * time.Millisecond,
	})
	return oracle, mux, func() {
		server.Close()
	}
}
