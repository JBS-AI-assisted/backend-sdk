// Package request decodes and validates JSON request bodies.
package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/JBS-AI-assisted/backend-sdk/model"
	"github.com/JBS-AI-assisted/backend-sdk/response"
)

// MaxBodyBytes caps every request body.
const MaxBodyBytes = 1 << 16

// validate is created once: the validator caches struct info and is safe for concurrent use.
var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Report fields by their JSON name ("newPassword"), as the contract does.
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// Decode reads r's JSON body into T and validates it with `validate` tags.
// On failure it writes 422 validation_failed and returns ok = false; the handler just returns.
func Decode[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var body T
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		response.WriteError(w, r, http.StatusUnprocessableEntity, model.CodeValidationFailed, "invalid JSON body")
		return body, false
	}
	if fields := Validate(body); len(fields) > 0 {
		response.WriteValidationError(w, r, fields)
		return body, false
	}
	return body, true
}

// Validate checks `validate` tags and returns contract field errors (nil when valid).
func Validate(v any) []model.FieldError {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return []model.FieldError{{Field: "", Code: model.FieldInvalidFormat}}
	}
	fields := make([]model.FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, model.FieldError{Field: fe.Field(), Code: fieldCode(fe)})
	}
	return fields
}

func fieldCode(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return model.FieldRequired
	case "min":
		if fe.Kind() == reflect.String {
			return model.FieldTooShort
		}
		return model.FieldOutOfRange
	case "max":
		if fe.Kind() == reflect.String {
			return model.FieldTooLong
		}
		return model.FieldOutOfRange
	case "gt", "gte", "lt", "lte", "oneof":
		return model.FieldOutOfRange
	default:
		return model.FieldInvalidFormat
	}
}
