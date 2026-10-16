// Package jwt issues and parses access tokens.
// Access tokens carry only the user ID (sub) and expiry: no email or name, since anyone holding a token can read it.
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Issuer is set on every access token and checked on parse.
const Issuer = "jsb-auth"

// ErrInvalidToken covers every reason a token is rejected (bad signature, expired, malformed).
var ErrInvalidToken = errors.New("invalid access token")

// Claims are the claims of an access token. UserID is in RegisteredClaims.Subject.
type Claims struct {
	jwt.RegisteredClaims
}

// GenerateAccessToken signs an HS256 access token for userID that expires after ttl.
func GenerateAccessToken(secret []byte, userID string, ttl time.Duration, now time.Time) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("jwt secret must be at least 32 bytes")
	}
	claims := Claims{RegisteredClaims: jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return signed, nil
}

// ParseAccessToken verifies signature, algorithm, issuer and expiry, and returns the claims.
func ParseAccessToken(secret []byte, token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(*jwt.Token) (any, error) { return secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(Issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("%w: missing sub", ErrInvalidToken)
	}
	return claims, nil
}
