// Package headlines handles querying and interacting with newsapi's top-headlines endpoint.
package headlines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/httperror"
	"github.com/riacataquian/news/internal/newsclient"
)

// PathPrefix is the newsapi's top headlines endpoint prefix.
const PathPrefix = "/top-headlines"

const maxPageSize = 100

const (
	// ErrMixParams is the error message for mixing parameters that shouldn't be mixed.
	// See Request Parameters > https://newsapi.org/docs/endpoints/top-headlines.
	ErrMixParams = "mixing `sources` with the `country` and `category` params"

	// ErrNoRequiredParams is the error message if no request parameter is present in the request.
	ErrNoRequiredParams = "required parameters are missing: sources, query, language, country, category."
)

// Endpoint is top headlines' request endpoint.
var Endpoint = newsclient.APIBaseURL + PathPrefix

// Params is the request parameters for top headlines endpoint.
// Requests should have at least one of these parameters.
// See Request Parameters > https://newsapi.org/docs/endpoints/top-headlines.
//
// It implements newsclient.Params interface.
type Params struct {
	// Country cannot be mixed with `sources` param.
	Country string `schema:"country"`
	// Category cannot be mixed with `sources` param.
	Category string `schema:"category"`
	// Sources is a comma-separated news sources.
	// See https://newsapi.org/sources for options.
	Sources string `schema:"sources"`
	// Query are keywords or phrase to search for.
	Query    string `schema:"query"`
	PageSize int    `schema:"pageSize"` // default: 20, maximum: 100
	Page     int    `schema:"page"`
}

// Client is an HTTP news API client.
// It implements the newsclient.Client interface.
type Client struct {
	newsclient.ServiceEndpoint
	ContextOrigin context.Context
	RequestOrigin *http.Request
}

// Get dispatches an HTTP GET request to the newsapi's top headlines endpoint.
// It times out after 5 seconds. // // It looks up for an env variable API_KEY and when found, set it to the request's header, // it then encodes params and set is as the request's query parameter.
//
// Finally, it dispatches the request by calling DispatchRequest method
// then encode the response accordingly.
func (c Client) Get(ctxOrigin context.Context, reqOrigin *http.Request, params newsclient.Params) (*news.Response, error) {
	ctx, cancel := context.WithTimeout(ctxOrigin, 5*time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("encoding query parameters: %v", err),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    newsclient.DocsBaseURL + "/endpoints" + PathPrefix,
		}
	}

	// Requests to external services should timeout for 5 seconds.
	req = req.WithContext(ctx)

	// Find and set request's API_KEY header.
	err = newsclient.LookupAndSetAuth(req)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    newsclient.DocsBaseURL + "/authentication",
		}
	}

	// Encode query parameters from the request origin.
	q, err := params.Encode()
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("encoding query parameters: %v", err),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    newsclient.DocsBaseURL + "/endpoints" + PathPrefix,
		}
	}
	req.URL.RawQuery = q

	// Dispatch request to news API.
	resp, err := c.DispatchRequest(req)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("fetching top headlines: %v", err),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    newsclient.DocsBaseURL + "/endpoints" + PathPrefix,
		}
	}

	return resp, nil
}

// DispatchRequest dispatches given r http.Request.
//
// It encodes and return a news.ErrorResponse when an error is encountered.
// Returns news.Response otherwise for successful requests.
func (c Client) DispatchRequest(r *http.Request) (*news.Response, error) {
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

// Encode encodes a p Params into a query string format. (e.g., foo=bar&wat=lol)
//
// It implements Params interface.
func (p Params) Encode() (string, error) {
	q := url.Values{}

	if p.Query != "" {
		q.Add("q", p.Query)
	}

	sources := p.Sources
	if sources != "" {
		q.Add("sources", sources)
	}

	if p.Country != "" {
		if sources != "" {
			return "", errors.New(ErrMixParams)
		}

		q.Add("country", p.Country)
	}

	if p.Category != "" {
		if sources != "" {
			return "", errors.New(ErrMixParams)
		}

		q.Add("category", p.Category)
	}

	// At this point, after all required parameters are evaluated and none is present,
	// return an ErrNoRequiredParams error.
	if q.Encode() == "" {
		return "", errors.New(ErrNoRequiredParams)
	}

	if p.Page != 0 {
		p := strconv.Itoa(p.Page)
		q.Add("page", p)
	}

	if p.PageSize != 0 {
		if p.PageSize > maxPageSize {
			return "", fmt.Errorf("the maximum page size is %d, you requested %d", maxPageSize, p.PageSize)
		}

		p := strconv.Itoa(p.PageSize)
		q.Add("pageSize", p)
	}

	return q.Encode(), nil // encodes q to bar=baz&foo=quux format.
}
