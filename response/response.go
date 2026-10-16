// Package response writes JSON bodies and contract-shaped errors.
package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/JBS-AI-assisted/backend-sdk/model"
	"github.com/JBS-AI-assisted/backend-sdk/requestid"
)

// WriteJSON writes body as JSON with the given status.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("write response", "err", err)
	}
}

// WriteError writes the shared Error body. The request ID comes from r's context.
func WriteError(w http.ResponseWriter, r *http.Request, status int, code model.ErrorCode, message string) {
	WriteJSON(w, status, model.Error{
		Code:      code,
		Message:   message,
		RequestID: requestid.FromContext(r.Context()),
	})
}

// WriteValidationError writes 422 validation_failed with per-field codes.
func WriteValidationError(w http.ResponseWriter, r *http.Request, fields []model.FieldError) {
	WriteJSON(w, http.StatusUnprocessableEntity, model.Error{
		Code:      model.CodeValidationFailed,
		Message:   "validation failed",
		Fields:    fields,
		RequestID: requestid.FromContext(r.Context()),
	})
}
