// Package newsclient provides functions, helpers and interface definition
// to interact with the external service https://newsapi.org.
package newsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/httperror"
)

// Params describes a Client's parameters.
type Params interface {
	Encode() (string, error)
}

// HTTPClient describes an HTTP client.
type HTTPClient interface {
	Get(context.Context, string, Params) (*news.Response, error)
}

// ServiceEndpoint wraps the URLs for newsapi endpoints.
type ServiceEndpoint struct {
	RequestURL string
	DocsURL    string
}

// Client performs HTTP requests to newsapi endpoints.
//
// See https://newsapi.org/docs/endpoints for the list of available endpoints.
type Client struct {
	ServiceEndpoint

	// Unexported fields.
	ctx context.Context
}

// NewFromContext returns a new Client.
// The supplied context.Context can be used later for HTTP request timeouts and cancellations.
func NewFromContext(ctx context.Context, se ServiceEndpoint) *Client {
	return &Client{ctx: ctx, ServiceEndpoint: se}
}

// Get fetches news from newsapi endpoints.
//
// Get injects the client's context, if any, to the current request.
// This can be used to enforce timeouts and cancellations.
//
// The request's `X-Api-Key` header is set with the supplied authKey.
func (client *Client) Get(ctx context.Context, authKey string, params Params) (*news.Response, error) {
	// Encode query parameters from the request origin.
	q, err := params.Encode()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, client.RequestURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Inject request authentication key.
	req.Header.Set("X-Api-Key", authKey)
	req.URL.RawQuery = q

	// Dispatch HTTP request to newsapi.
	resp, err := dispatchReq(req)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("error while dispatching request: %v", err),
			DocsURL: client.DocsURL,
		}
	}
	return resp, nil
}

// dispatchReq dispatches the supplied http.Request.
//
// It encodes and return a news.ErrorResponse when an error is encountered.
// Returns news.Response otherwise for successful requests.
func dispatchReq(r *http.Request) (*news.Response, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var res news.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
		return nil, &res
	}

	var res news.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &res, nil
}
