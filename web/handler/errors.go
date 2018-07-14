package handler

// This file contains error handlers.

import (
	"context"
	"net/http"

	"github.com/riacataquian/news/pkg/httperror"
)

// NotFound handles HTTP requests for missing or not found pages and resources.
func NotFound(_ context.Context) ([]byte, error) {
	return nil, httperror.New(http.StatusNotFound, "page not found")
}
