package newsclient

import "time"

// This file handles querying and interacting with newsapi's everthing endpoint.

// EverythingPathPrefix is newsapi's everything endpoint prefix.
const EverythingPathPrefix = "/everything"

// Sorting is the order to sort articles in.
type Sorting string

const (
	// relevancy means articles more closely related to Query comes first.
	relevancy Sorting = "relevancy"
	// popularity means articles from popular sources and publishers comes first.
	popularity Sorting = "popularity"
	// publishedAt means newest articles comes first.
	publishedAt Sorting = "publishedAt"
)

// Language is a 2-letter IS0-639-1 code of the language to get the news for.
type Language string

const (
	ar Language = "ar"
	de Language = "de"
	en Language = "en"
	es Language = "es"
	fr Language = "fr"
	he Language = "he"
	it Language = "it"
	nl Language = "nl"
	no Language = "no"
	pt Language = "pt"
	ru Language = "ru"
	se Language = "se"
	ud Language = "ud"
	zh Language = "zh"
)

// EverythingParams is the request parameters for news under everything category.
// All of which are optional parameters, except the `apiKey`,
// which in newsclient's case is sent as a request header.
// See Request Parameters > https://newsapi.org/docs/endpoints/everything.
//
// It implements Params interface.
type EverythingParams struct {
	// Query are keywords or phrase to search for.
	Query string `schema:"query"`
	// Sources is a comma-separated news sources.
	Sources string `schema:"sources"`
	// Domains are comma-separated string of domains to restrict the search to.
	Domains string `schema:"domains"`
	// To is the date and optional time for the newest article allowed.
	// Expects an ISO format, i.e., 2018-07-28 or 2018-07-28T14:28:41.
	To time.Time `schema:"To"`
	// From is the date and optional time for the oldest article allowed.
	// Expects an ISO format, i.e., 2018-07-28 or 2018-07-28T14:28:41.
	From     time.Time           `schema:"from"`
	Language `schema:"language"` // defaults to all languages returned.
	SortBy   Sorting             `schema:"sortBy"`   // defaults to publishedAt.
	PageSize int                 `schema:"pageSize"` // default: 20, maximum: 100
	Page     int                 `schema:"page"`
}
