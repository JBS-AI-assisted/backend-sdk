// Package model holds the types shared by every service. It mirrors
// backend/api/common.yaml: change the contract first, then this file.
package model

// ErrorCode is the stable, machine-readable error code. The app maps on it, never on Message.
type ErrorCode string

const (
	CodeValidationFailed    ErrorCode = "validation_failed"
	CodeUnauthorized        ErrorCode = "unauthorized"
	CodeInvalidCredentials  ErrorCode = "invalid_credentials"
	CodeEmailTaken          ErrorCode = "email_taken"
	CodeInvalidRefreshToken ErrorCode = "invalid_refresh_token"
	CodeInvalidCode         ErrorCode = "invalid_code"
	CodeCodeLocked          ErrorCode = "code_locked"
	CodeInvalidPassword     ErrorCode = "invalid_password"
	CodeProfileNotFound     ErrorCode = "profile_not_found"
	CodeRateLimited         ErrorCode = "rate_limited"
	CodeInternalError       ErrorCode = "internal_error"
)

// Field-level codes used in FieldError.Code.
const (
	FieldRequired      = "required"
	FieldInvalidFormat = "invalid_format"
	FieldTooShort      = "too_short"
	FieldTooLong       = "too_long"
	FieldOutOfRange    = "out_of_range"
	FieldTooCommon     = "too_common"
	FieldUnderMinAge   = "under_min_age"
)

// FieldError is one per-field problem in a validation_failed response.
type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code"`
}

// Error is the body of every non-2xx response.
type Error struct {
	Code      ErrorCode    `json:"code"`
	Message   string       `json:"message"`
	Fields    []FieldError `json:"fields,omitempty"`
	RequestID string       `json:"requestId,omitempty"`
}
