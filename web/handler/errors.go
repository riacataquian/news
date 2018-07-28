package handler

// This file contains error handlers.

import (
	"context"
	"net/http"

	"github.com/riacataquian/news/internal/httperror"
)

// NotFound handles HTTP requests for missing or not found pages and resources.
func NotFound(_ context.Context, r *http.Request) (*SuccessResponse, error) {
	return nil, &httperror.HTTPError{
		Code:       http.StatusNotFound,
		Message:    "page not found",
		RequestURL: r.URL.String(),
	}
}
