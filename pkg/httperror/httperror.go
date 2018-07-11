// Package httperror contains functions and helpers for handling HTTP errors.
package httperror

import (
	"fmt"
)

// HTTPError describes an HTTP error.
type HTTPError struct {
	Code          int `json:"statusCode"`
	ErrorResponse `json:"error"`
}

// ErrorResponse is a generic error object.
type ErrorResponse struct {
	Message string     `json:"message"`
	Errors  []FieldErr `json:"errors,omitempty"`
}

// FieldErr describes an error for a resource or field.
type FieldErr struct {
	Field  string
	Errors []string
}

// New returns a new HTTPError with the supplied code, message and errors.
func New(code int, msg string, errs ...FieldErr) *HTTPError {
	return &HTTPError{code, ErrorResponse{msg, errs}}
}

// Error formats and return an HTTPError's message.
func (h *HTTPError) Error() string {
	if h.Errors == nil {
		return h.Message
	}

	return fmt.Sprintf(`%s See "Errors" field for more info.`, h.Errors)
}
