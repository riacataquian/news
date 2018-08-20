// Package newsclient provides functions, helpers and interface definition
// to interact with the external service https://newsapi.org.
package newsclient

import (
	"context"
	"net/http"

	"github.com/riacataquian/news/api/news"
)

var (
	// APIBaseURL is the base URL of newsapi's API endpoint.
	APIBaseURL = "https://newsapi.org/v2"

	// DocsBaseURL is the base URL of newsapi's documentation.
	DocsBaseURL = "https://newsapi.org/docs"
)

// Params describes a Client's parameters.
type Params interface {
	Encode() (string, error)
}

// Client describes an HTTP news client.
type Client interface {
	NewGetRequest(context.Context) (*http.Request, error)
	AuthorizeReq(*http.Request, string)
	Get(*http.Request, Params) (*news.Response, error)
}
