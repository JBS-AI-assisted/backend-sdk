package request

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/JBS-AI-assisted/backend-sdk/model"
)

type signUp struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=10,max=128"`
	Year     int    `json:"birthYear" validate:"omitempty,min=1900"`
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantOK     bool
		wantFields []model.FieldError
	}{
		{"valid", `{"email":"a@b.co","password":"long-enough-pw"}`, true, nil},
		{"bad json", `{"email":`, false, nil},
		{"unknown field", `{"email":"a@b.co","password":"long-enough-pw","admin":true}`, false, nil},
		{"missing email", `{"password":"long-enough-pw"}`, false, []model.FieldError{{Field: "email", Code: model.FieldRequired}}},
		{"bad email", `{"email":"nope","password":"long-enough-pw"}`, false, []model.FieldError{{Field: "email", Code: model.FieldInvalidFormat}}},
		{"short password", `{"email":"a@b.co","password":"short"}`, false, []model.FieldError{{Field: "password", Code: model.FieldTooShort}}},
		{"year out of range", `{"email":"a@b.co","password":"long-enough-pw","birthYear":1800}`, false, []model.FieldError{{Field: "birthYear", Code: model.FieldOutOfRange}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			_, ok := Decode[signUp](rec, req)

			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v (body %s)", ok, tt.wantOK, rec.Body.String())
			}
			if ok {
				return
			}
			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422", rec.Code)
			}
			var got model.Error
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if got.Code != model.CodeValidationFailed {
				t.Errorf("code = %q, want validation_failed", got.Code)
			}
			if tt.wantFields != nil && !reflect.DeepEqual(got.Fields, tt.wantFields) {
				t.Errorf("fields = %+v, want %+v", got.Fields, tt.wantFields)
			}
		})
	}
}
