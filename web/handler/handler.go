// Package handler contains routes and HTTP handlers.
package handler

import (
	"context"
	"net/http"

	"github.com/riacataquian/news/internal/store"
)

// SuccessResponse describes a successful HTTP response.
type SuccessResponse struct {
	Code       int    `json:"code"`
	RequestURL string `json:"requestURL"`
	// Count is the queried result count.
	Count int `json:"count"`
	// Page is the current result's page.
	Page int `json:"page"`
	// TotalCount is the number of queryable results.
	TotalCount int `json:"totalCount"`
	// Data is the actual response from newsapi.
	Data interface{} `json:"data"`
}

// Func describes a function that handles HTTP requests and responses.
type Func func(context.Context, store.Store, *http.Request) (*SuccessResponse, error)

// Routes is the lookup table for URL paths and their matching handlers.
var Routes = []struct {
	Path        string
	HandlerFunc Func
}{
	{"/list", List},
	{"/headlines", TopHeadlines},
	{"/{*}", NotFound},
}
