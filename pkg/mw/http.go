package mw

import (
	"fmt"
	"net/http"
	"log"
)

// Adds json content type to a handler
func AddJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// Adds request logging to a handler
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s:\t%s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
