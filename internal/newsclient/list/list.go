// Package list contains constants, endpoints, params and errors for news list.
package list

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/riacataquian/news/internal/newsclient"
)

var (
	// ErrNoRequiredParams is the error message if no request parameter is present in the request.
	ErrNoRequiredParams = errors.New("required parameters are missing: query, sources, domains.")

	// ServiceEndpoint wraps URLs to newsapi's everything endpoint.
	ServiceEndpoint = newsclient.ServiceEndpoint{
		RequestURL: "https://newsapi.org/v2/everything",
		DocsURL:    "https://newsapi.org/docs/endpoints/everything",
	}
)

// Sorting is the order to sort articles in.
type Sorting string

const (
	// Relevancy means articles more closely related to Query comes first.
	Relevancy Sorting = "relevancy"
	// Popularity means articles from popular sources and publishers comes first.
	Popularity Sorting = "popularity"
	// PublishedAt means newest articles comes first.
	PublishedAt Sorting = "publishedAt"
	// maxPageSize is the maximum page size for requesting list news.
	maxPageSize = 100
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

// Encode encodes an list's Params into a query string format. (e.g., foo=bar&wat=lol)
//
// It implements newsclient.Params interface.
func (p *Params) Encode() (string, error) {
	if p == nil {
		return "", ErrNoRequiredParams
	}

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
		return "", ErrNoRequiredParams
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
