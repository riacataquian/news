// Package list handles querying and interacting with newsapi's everything endpoint.
package list

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/httpclient"
	"github.com/riacataquian/news/internal/httperror"
	"github.com/riacataquian/news/internal/newsclient"
)

// PathPrefix is newsapi's everything endpoint prefix.
const PathPrefix = "/everything"

const maxPageSize = 100

// ErrNoRequiredParams is the error message if no request parameter is present in the request.
const ErrNoRequiredParams = "required parameters are missing: query, sources, domains."

// ErrMissingAPIKey is the error message for missing API key.
var ErrMissingAPIKey = errors.New("missing API key in the request header or parameters")

// Endpoint is the request endpoint.
var Endpoint = newsclient.APIBaseURL + PathPrefix

var client = httpclient.NewClient()

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

// Params is the request parameters for list news request.
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
	From string `schema:"from"`
	// To is the date and optional time for the newest article allowed.
	// Expects an ISO format, i.e., 2018-07-28 or 2018-07-28T14:28:41.
	To string `schema:"To"`
	// Language is a 2-letter IS0-639-1 code of the language to get the news for.
	// See Request Parameters > language > https://newsapi.org/docs/endpoints/everything.
	Language string  `schema:"language"` // defaults to all languages returned.
	SortBy   Sorting `schema:"sortBy"`   // defaults to publishedAt.
	PageSize int     `schema:"pageSize"` // default: 20, maximum: 100
	Page     int     `schema:"page"`
}

// Client is an HTTP newsapi client.
//
// Client satisfies the newsclient.Client interface.
type Client struct {
	serviceEndpoint
}

// serviceEndpoint wraps domains and their external URLs - where requests should be dispatched to.
type serviceEndpoint struct {
	everything string
}

// NewClient returns a new list.Client.
func NewClient() newsclient.Client {
	return &Client{
		serviceEndpoint: serviceEndpoint{
			everything: Endpoint,
		},
	}
}

// NewGetRequest ...
func (c *Client) NewGetRequest(ctx context.Context) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, c.everything, nil)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)
	return r, nil
}

// AuthorizeReq ...
func (c *Client) AuthorizeReq(r *http.Request, key string) {
	r.Header.Set("X-Api-Key", key)
}

// Get ...
func (c *Client) Get(r *http.Request, params newsclient.Params) (*news.Response, error) {
	if k := r.Header.Get("X-Api-Key"); k == "" {
		return nil, ErrMissingAPIKey
	}

	// Encode query parameters from the request origin.
	q, err := params.Encode()
	if err != nil {
		return nil, err
	}
	r.URL.RawQuery = q

	// Dispatch request to newsapi.
	resp, err := client.DispatchRequest(r)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("fetching news: %v", err),
			DocsURL: newsclient.DocsBaseURL + "/endpoints" + PathPrefix,
		}
	}

	return resp, nil
}

// Encode encodes a p Params into a query string format. (e.g., foo=bar&wat=lol)
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
	if from != "" {
		q.Add("from", from)
	}

	to := p.To
	if to != "" {
		q.Add("to", to)
	}

	language := p.Language
	if language != "" {
		q.Add("language", language)
	}

	sortBy := p.SortBy
	if sortBy != "" {
		q.Add("sortBy", string(sortBy))
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
