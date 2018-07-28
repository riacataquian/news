// Package everything handles querying and interacting with newsapi's everything endpoint.
package everything

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

// This file handles querying and interacting with newsapi's everthing endpoint.

// EverythingPathPrefix is newsapi's everything endpoint prefix.
const EverythingPathPrefix = "/everything"

// Sorting is the order to sort articles in.
type Sorting string

const (
	// Relevancy means articles more closely related to Query comes first.
	Relevancy Sorting = "relevancy"
	// Popularity means articles from popular sources and publishers comes first.
	Popularity Sorting = "popularity"
	// PublishedAt means newest articles comes first.
	PublishedAt Sorting = "publishedAt"
)

// Language is a 2-letter IS0-639-1 code of the language to get the news for.
type Language string

const (
	AR Language = "ar"
	DE Language = "de"
	EN Language = "en"
	ES Language = "es"
	FR Language = "fr"
	HE Language = "he"
	IT Language = "it"
	NL Language = "nl"
	NO Language = "no"
	PT Language = "pt"
	RU Language = "ru"
	SE Language = "se"
	UD Language = "ud"
	ZH Language = "zh"
)

// Params is the request parameters for news under everything category.
// Requests should have at least one of these parameters.
// See Request Parameters > https://newsapi.org/docs/endpoints/everything.
//
// It implements newsclient.Params interface.
type Params struct {
	// Query are keywords or phrase to search for.
	Query string `schema:"query"`
	// Sources is a comma-separated news sources.
	// See https://newsapi.org/sources for options.
	Sources string `schema:"sources"`
	// Domains are comma-separated string of domains to restrict the search to.
	Domains string `schema:"domains"`
	// From is the date and optional time for the oldest article allowed.
	// Expects an ISO format, i.e., 2018-07-28 or 2018-07-28T14:28:41.
	From *time.Time `schema:"from"`
	// To is the date and optional time for the newest article allowed.
	// Expects an ISO format, i.e., 2018-07-28 or 2018-07-28T14:28:41.
	To       *time.Time          `schema:"To"`
	Language `schema:"language"` // defaults to all languages returned.
	SortBy   Sorting             `schema:"sortBy"`   // defaults to publishedAt.
	PageSize int                 `schema:"pageSize"` // default: 20, maximum: 100
	Page     int                 `schema:"page"`
}

// Client is an HTTP news API client.
// It implements the newsclient.Client interface.
type Client struct {
	newsclient.ServiceEndpoint
	ContextOrigin context.Context
	RequestOrigin *http.Request
}

// Get dispatches an HTTP GET request to the newsapi's everthing endpoint.
// It times out after 5 seconds.
//
// It looks up for an env variable API_KEY and when found, set it to the request's header,
// it then encodes params and set is as the request's query parameter.
//
// Finally, it dispatches the request by calling DispatchRequest method
// then encode the response accordingly.
func (c Client) Get(ctxOrigin context.Context, reqOrigin *http.Request, params Params) (*news.Response, error) {
	ctx, cancel := context.WithTimeout(ctxOrigin, 5*time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, c.URL, nil)
	if err != nil {
		return nil, err
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
			DocsURL:    newsclient.DocsBaseURL + "/endpoints" + EverythingPathPrefix,
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
			DocsURL:    newsclient.DocsBaseURL + "/endpoints" + EverythingPathPrefix,
		}
	}

	return resp, err
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

// Encode encodes a p EverythingParams into a query string format. (e.g., foo=bar&wat=lol)
//
// It implements newsclient.Params interface.
func (p Params) Encode() (string, error) {
	q := url.Values{}

	if p.Query != "" {
		q.Add("q", p.Query)
	}

	sources := p.Sources
	if sources != "" {
		q.Add("sources", sources)
	}

	domains := p.Domains
	if domains != "" {
		q.Add("domains", domains)
	}

	from := p.From
	if from != nil {
		q.Add("from", from.Format(time.RFC3339))
	}

	to := p.To
	if to != nil {
		q.Add("to", to.Format(time.RFC3339))
	}

	language := p.Language
	if language != "" {
		q.Add("language", string(language))
	}

	sortBy := p.SortBy
	if sortBy != "" {
		q.Add("sortBy", string(sortBy))
	}

	// At this point, after all required parameters are evaluated and none is present,
	// return an ErrNoRequiredParams error.
	if q.Encode() == "" {
		return "", errors.New(newsclient.ErrNoRequiredParams)
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
