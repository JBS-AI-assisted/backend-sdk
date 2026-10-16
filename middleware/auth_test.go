package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	sdkjwt "github.com/JBS-AI-assisted/backend-sdk/jwt"
)

var secret = []byte("test-secret-test-secret-test-secret!")

func TestAuth(t *testing.T) {
	now := time.Now()
	valid, err := sdkjwt.GenerateAccessToken(secret, "user-1", 15*time.Minute, now)
	if err != nil {
		t.Fatal(err)
	}
	expired, _ := sdkjwt.GenerateAccessToken(secret, "user-1", time.Minute, now.Add(-time.Hour))
	otherKey, _ := sdkjwt.GenerateAccessToken([]byte("another-secret-another-secret-12345"), "user-1", time.Minute, now)
	noneAlg, _ := gojwt.NewWithClaims(gojwt.SigningMethodNone, gojwt.MapClaims{
		"sub": "user-1", "iss": sdkjwt.Issuer, "exp": now.Add(time.Hour).Unix(),
	}).SignedString(gojwt.UnsafeAllowNoneSignatureType)
	noExp, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.MapClaims{
		"sub": "user-1", "iss": sdkjwt.Issuer,
	}).SignedString(secret)
	wrongTypes, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.MapClaims{
		"sub": 42, "iss": sdkjwt.Issuer, "exp": "tomorrow",
	}).SignedString(secret)

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"valid", "Bearer " + valid, http.StatusOK},
		{"no header", "", http.StatusUnauthorized},
		{"no bearer prefix", valid, http.StatusUnauthorized},
		{"empty token", "Bearer ", http.StatusUnauthorized},
		{"garbage", "Bearer abc.def.ghi", http.StatusUnauthorized},
		{"expired", "Bearer " + expired, http.StatusUnauthorized},
		{"other key", "Bearer " + otherKey, http.StatusUnauthorized},
		{"alg none", "Bearer " + noneAlg, http.StatusUnauthorized},
		{"no exp", "Bearer " + noExp, http.StatusUnauthorized},
		{"wrong claim types (used to panic)", "Bearer " + wrongTypes, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotID string
			h := Auth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotID, _ = UserID(r.Context())
			}))
			req := httptest.NewRequest(http.MethodGet, "/v1/me", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK && gotID != "user-1" {
				t.Errorf("user id = %q, want user-1", gotID)
			}
			if tt.wantStatus == http.StatusUnauthorized && !strings.Contains(rec.Body.String(), `"code":"unauthorized"`) {
				t.Errorf("body = %s, want code unauthorized", rec.Body.String())
			}
		})
	}
}

func TestGenerateAccessTokenRejectsShortSecret(t *testing.T) {
	if _, err := sdkjwt.GenerateAccessToken([]byte("short"), "user-1", time.Minute, time.Now()); err == nil {
		t.Fatal("want error for a short secret")
	}
}
