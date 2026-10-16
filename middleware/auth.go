package middleware

import (
	"context"
	"net/http"
	"strings"

	sdkjwt "github.com/JBS-AI-assisted/backend-sdk/jwt"
	"github.com/JBS-AI-assisted/backend-sdk/model"
	"github.com/JBS-AI-assisted/backend-sdk/response"
)

type userIDKey struct{}

// Auth rejects requests without a valid "Authorization: Bearer <access token>" with 401 unauthorized,
// and puts the user ID into the request context (read it with UserID).
func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if !ok || token == "" {
				unauthorized(w, r)
				return
			}
			claims, err := sdkjwt.ParseAccessToken(secret, token)
			if err != nil {
				unauthorized(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey{}, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserID returns the authenticated user's ID. ok is false outside Auth.
func UserID(ctx context.Context) (id string, ok bool) {
	id, ok = ctx.Value(userIDKey{}).(string)
	return id, ok && id != ""
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	response.WriteError(w, r, http.StatusUnauthorized, model.CodeUnauthorized, "missing, invalid or expired access token")
}
