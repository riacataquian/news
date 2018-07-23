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

	// ErrMixParams is the error message for mixing parameters that shouldn't be mixed.
	// See Request Parameters > https://newsapi.org/docs/endpoints/top-headlines.
	ErrMixParams = "mixing `sources` with the `country` and `category` params"

	// ErrNoRequiredParams is the error message if no parameter is present in the request.
	ErrNoRequiredParams = "required parameters are missing: sources, q, language, country, category."
)

// Params describes a Client's parameters.
type Params interface {
	Encode() (string, error)
}

// Client describes an HTTP news client.
type Client interface {
	TopHeadlines(context.Context, *http.Request, Params) (*news.Response, error)
	DispatchRequest(*http.Request) (*news.Response, error)
}

// ServiceEndpoint wraps a URL in where a request should be dispatched to.
type ServiceEndpoint struct {
	URL string
}

// NewsClient is an HTTP news API client.
// It implements the Client interface.
type NewsClient struct {
	ServiceEndpoint
	ContextOrigin context.Context
	RequestOrigin *http.Request
}

// lookupAndSetAuth ...
func lookupAndSetAuth(r *http.Request) error {
	k, ok := os.LookupEnv("API_KEY")
	if !ok {
		return errors.New(ErrMissingAPIKey)
	}

	r.Header.Set("X-Api-Key", k)
	return nil
}
