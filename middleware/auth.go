package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"sdk/model"
	"sdk/response"
	"strings"

	JWThelper "sdk/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
	
		if authorizationHeader == "" {
			response.WriteJSON(w, model.ErrorResponse{
				Message: model.ErrorUnauthorized.Error(),
			}, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			response.WriteJSON(w, model.ErrorResponse{
				Message: model.ErrorUnauthorized.Error(),
			}, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		
		if err != nil || !token.Valid {
			response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorUnauthorized.Error()}, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorUnauthorized.Error()}, http.StatusUnauthorized)
			return
		}

		if claims["exp"] == nil || int64(claims["exp"].(float64)) < time.Now().Unix() {
			response.WriteJSON(w, model.ErrorResponse{Message: model.ErrorUnauthorized.Error()}, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), JWThelper.UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	}
}

func ExtractClaims(ctx context.Context) (JWThelper.JWTDataClaims, bool) {
	claims, ok := ctx.Value(JWThelper.UserClaimsKey).(jwt.MapClaims)

	if !ok {
		return JWThelper.JWTDataClaims{}, false
	}

	idFloat, _ := claims["id"].(float64)

	jwtClaims := JWThelper.JWTDataClaims{
		Email: claims["email"].(string),
		ID:    uint(idFloat),
		Name:  claims["name"].(string),
		Exp:   int64(claims["exp"].(float64)),
	}

	return jwtClaims, true
}