package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/JBS-AI-assisted/backend-sdk/requestid"
)

// Logger logs one line per request: method, path, status, duration, request ID.
// It never logs bodies, query strings, headers or client IPs (personal data stays out of logs).
// Put it after requestid.Middleware in the chain.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", requestid.FromContext(r.Context()),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Unwrap lets http.ResponseController reach the underlying writer (flush for SSE).
func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }
