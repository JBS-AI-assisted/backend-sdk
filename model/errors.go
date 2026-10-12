package model

import "errors"

var (
	ErrorMissingToken = errors.New("Error missing token")
	ErrorInvalidToken = errors.New("Error invalid token")
	ErrorUnauthorized = errors.New("Error unauthorized")
)

type ErrorResponse struct {
	Message string      `json:"message"`
	TraceId string      `json:"traceId"`
	Payload interface{} `json:"payload"`
}