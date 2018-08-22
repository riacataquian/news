// Package httperror contains functions and helpers for handling HTTP errors.
package httperror

import (
	"fmt"
)

// HTTPError describes an HTTP error.
type HTTPError struct {
	Code        int           `json:"statusCode"`
	Message     string        `json:"message"`
	RequestURL  string        `json:"requestUrl,omitempty"`
	DocsURL     string        `json:"docsUrl,omitempty"`
	FieldErrors []FieldErrors `json:"errors,omitempty"`
}

// FieldErrors is a generic error object.
type FieldErrors struct {
	Message string     `json:"message"`
	Errors  []FieldErr `json:"errors,omitempty"`
}

// FieldErr describes an error for a resource or field.
type FieldErr struct {
	Field  string
	Errors []string
}

// Error formats and return an HTTPError's message.
//
// Error is the HTTPError's error implementation.
func (e *HTTPError) Error() string {
	if len(e.FieldErrors) == 0 {
		return e.Message
	}

	return fmt.Sprintf(`%s. See "errors" field for more info.`, e.Message)
}
