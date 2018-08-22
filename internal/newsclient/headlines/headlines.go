// Package headlines contains constants, endpoints, params and errors for top headlines.
package headlines

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/riacataquian/news/internal/newsclient"
)

var (
	// ErrMixParams is the error message for mixing parameters that shouldn't be mixed.
	// See Request Parameters > https://newsapi.org/docs/endpoints/top-headlines.
	ErrMixParams = errors.New("mixing `sources` with the `country` and `category` params")

	// ErrNoRequiredParams is the error message if no request parameter is present in the request.
	ErrNoRequiredParams = errors.New("required parameters are missing: sources, query, language, country, category.")

	// ErrInvalidPageSize is the error message if the supplied maximum page size exceeded the allowed size.
	ErrInvalidPageSize = errors.New("invalid page, maximum page size is 100")

	// ServiceEndpoint wraps URLs to newsapi's top-headlines endpoint.
	ServiceEndpoint = newsclient.ServiceEndpoint{
		RequestURL: "https://newsapi.org/v2/top-headlines",
		DocsURL:    "https://newsapi.org/docs/endpoints/top-headlines",
	}
)

// maxPageSize is the maximum page size for requesting top headlines news.
const maxPageSize = 100

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

// Encode encodes a headlines' Params into a query string format. (e.g., foo=bar&wat=lol)
//
// It implements Params interface.
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

	if p.Country != "" {
		if sources != "" {
			return "", ErrMixParams
		}

		q.Add("country", p.Country)
	}

	if p.Category != "" {
		if sources != "" {
			return "", ErrMixParams
		}

		q.Add("category", p.Category)
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
			return "", ErrInvalidPageSize
		}

		p := strconv.Itoa(p.PageSize)
		q.Add("pageSize", p)
	}

	return q.Encode(), nil // encodes q to bar=baz&foo=quux format.
}
