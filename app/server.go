package main

import (
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Access logs
	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			log.Printf("[%s] %s, %s\n", r.Method, r.RemoteAddr, r.URL.Path)
		})
	}(s.mux)

	// Panic resolver
	handler = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal server error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}(handler)

	// JSON response
	handler = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			next.ServeHTTP(w, r)
		})
	}(handler)

	handler.ServeHTTP(w, r)
}

func NewServer() *Server {
	return &Server{mux: new(http.ServeMux)}
}
