// Package handler contains routes and HTTP handlers.
package handler

import (
	"context"
	"net/http"
)

// SuccessResponse describes a sucessful HTTP response.
type SuccessResponse struct {
	Code       int         `json:"code"`
	RequestURL string      `json:"requestURL"`
	Data       interface{} `json:"data"`
}

// Func describes a function that handles HTTP requests and responses.
type Func func(context.Context, *http.Request) (*SuccessResponse, error)

// Routes is the lookup table for URL paths and their matching handlers.
var Routes = []struct {
	Path        string
	HandlerFunc Func
}{
	{"/headlines", News},
	{"/{*}", NotFound},
}
