// Package handler contains routes and HTTP handlers.
package handler

import (
	"context"
)

// Func describes a function that handles HTTP requests and responses.
type Func func(context.Context) ([]byte, error)

// Routes is the lookup table for URL paths and their matching handlers.
var Routes = []struct {
	Path        string
	HandlerFunc Func
}{
	{"/{*}", NotFound},
}
