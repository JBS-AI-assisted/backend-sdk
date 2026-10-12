package middleware

import (
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		println("Incoming request:", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		println("Request processed in:", duration.String())
	})
}