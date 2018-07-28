// Package newsclient provides functions, helpers and interface definition
// to interact with the external service https://newsapi.org.
package newsclient

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/riacataquian/news/api/news"
)

var (
	// APIBaseURL is the base URL of newsapi's API endpoint.
	APIBaseURL = "https://newsapi.org/v2"

	// DocsBaseURL is the base URL of newsapi's documentation.
	DocsBaseURL = "https://newsapi.org/docs"

	// ErrMissingAPIKey is the error message for missing API key.
	ErrMissingAPIKey = "missing API key in the request header or parameters"
)

// Params describes a Client's parameters.
type Params interface {
	Encode() (string, error)
}

// Client describes an HTTP news client.
type Client interface {
	Get(context.Context, *http.Request, Params) (*news.Response, error)
	DispatchRequest(*http.Request) (*news.Response, error)
}

// ServiceEndpoint wraps a URL in where a request should be dispatched to.
type ServiceEndpoint struct {
	URL string
}

// LookupAndSetAuth sets the env variable API_KEY in the supplied request.
func LookupAndSetAuth(r *http.Request) error {
	k, ok := os.LookupEnv("API_KEY")
	if !ok {
		return errors.New(ErrMissingAPIKey)
	}

	r.Header.Set("X-Api-Key", k)
	return nil
}
