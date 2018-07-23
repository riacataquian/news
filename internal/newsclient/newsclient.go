// Package newsclient ...
package newsclient

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/riacataquian/news/api/news"
)

var (
	// APIBaseURL ...
	APIBaseURL = "https://newsapi.org/v2"

	// DocsBaseURL ...
	DocsBaseURL = "https://newsapi.org/docs"

	// ErrMissingAPIKey ...
	ErrMissingAPIKey = "missing API key in the request header"

	// ErrMixParams ...
	ErrMixParams = "mixing `sources` with the `country` and `category` params"

	// ErrNoRequiredParams
	ErrNoRequiredParams = "required parameters are missing: sources, q, language, country, category."
)

// Params ...
type Params interface {
	Encode() (string, error)
}

// Client ...
type Client interface {
	GetTopHeadlines(Params) (*news.Response, error)
	GetContextOrigin() context.Context
	GetRequestOrigin() *http.Request
	GetServiceEndpoint() ServiceEndpoint
	DispatchRequest(*http.Request) (*news.Response, error)
}

// ServiceEndpoint ...
type ServiceEndpoint struct {
	URL string
}

// NewsClient ...
//
// It implements the Client interface.
type NewsClient struct {
	ServiceEndpoint
	ContextOrigin context.Context
	RequestOrigin *http.Request
}

// GetContextOrigin ...
func (nc NewsClient) GetContextOrigin() context.Context {
	return nc.ContextOrigin
}

// GetRequestOrigin ...
func (nc NewsClient) GetRequestOrigin() *http.Request {
	return nc.RequestOrigin
}

// GetServiceEndpoint ...
func (nc NewsClient) GetServiceEndpoint() ServiceEndpoint {
	return nc.ServiceEndpoint
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
