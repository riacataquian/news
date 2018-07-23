package newsclient

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
)

// This file handles top headlines news.

// TopHeadlinesPathPrefix ...
var TopHeadlinesPathPrefix = "/top-headlines"

// TopHeadlinesParams ...
//
// All of which are optional parameters.
// It implements Params interface.
type TopHeadlinesParams struct {
	// It cannot be mixed with `sources` param.
	Country string `schema:"country"`
	// It cannot be mixed with `sources` param.
	Category string `schema:"category"`
	// Sources is a comma-separated sources.
	Sources  string `schema:"sources"`
	Query    string `schema:"query"`
	PageSize int    `schema:"pageSize"` // default: 20, maximum: 100
	Page     int    `schema:"page"`
}

// TopHeadlines ...
func (c NewsClient) TopHeadlines(ctxOrigin context.Context, reqOrigin *http.Request, params Params) (*news.Response, error) {
	ctx, cancel := context.WithTimeout(ctxOrigin, 5*time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, err
	}

	// Requests to external services should timeout for 5 seconds.
	req = req.WithContext(ctx)

	// Encode query parameters from the request origin.
	q, err := params.Encode()
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("encoding query parameters: %v", err),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    DocsBaseURL + "/endpoints" + TopHeadlinesPathPrefix,
		}
	}
	req.URL.RawQuery = q

	// Find and set request's API_KEY header.
	err = lookupAndSetAuth(req)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    DocsBaseURL + "/authentication",
		}
	}

	// Dispatch request to news API.
	resp, err := c.DispatchRequest(req)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    fmt.Sprintf("fetching top headlines: %v", err),
			RequestURL: reqOrigin.URL.String(),
			DocsURL:    DocsBaseURL + "/endpoints" + TopHeadlinesPathPrefix,
		}
	}

	return resp, err
}

// DispatchRequest dispatches given r http.Request.
//
// It encodes and return a news.ErrorResponse when an error is encountered.
// Returns news.Response otherwise for successful requests.
func (c NewsClient) DispatchRequest(r *http.Request) (*news.Response, error) {
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

// Encode encodes a p TopHeadlinesParams into a query string format. (e.g., foo=bar&wat=lol)
//
// It implements Params interface.
func (p TopHeadlinesParams) Encode() (string, error) {
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
		p := strconv.Itoa(p.PageSize)
		q.Add("pageSize", p)
	}

	return q.Encode(), nil // encodes q to bar=baz&foo=quux format.
}
