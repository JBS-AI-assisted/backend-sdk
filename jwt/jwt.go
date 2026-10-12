package jwt

import (
	"context"
	"fmt"
	"net/http"
	"sdk/model"
	"sdk/response"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTDataClaims struct {
	Email string
	ID uint
	Name string
	Exp int64
}

type contextKey string
const UserClaimsKey contextKey = "user_claims"

func JWTMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorMissingToken.Error()}, http.StatusUnauthorized)
				return
			}

			tokenString := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorInvalidToken.Error()}, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorInvalidToken.Error()}, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func GenerateToken(secretKey string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
