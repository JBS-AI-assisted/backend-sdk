// Package requestid gives every request an ID that is returned in errors and written to logs.
package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"regexp"
)

// Header is read from the caller (gateway, other service) and always set on the response.
const Header = "X-Request-ID"

type contextKey struct{}

// valid limits what we accept from outside, so a caller cannot inject junk into our logs.
var valid = regexp.MustCompile(`^[A-Za-z0-9-]{8,64}$`)

// Middleware reuses a valid incoming X-Request-ID or creates a new one.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(Header)
		if !valid.MatchString(id) {
			id = newID()
		}
		w.Header().Set(Header, id)
		next.ServeHTTP(w, r.WithContext(WithID(r.Context(), id)))
	})
}

// WithID stores id in ctx.
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext returns the request ID, or "" when there is none.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b) // crypto/rand.Read never returns an error since Go 1.24.
	return hex.EncodeToString(b)
}
