package cron

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/auth"
	"github.com/riacataquian/news/internal/clock"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/list"
	"github.com/riacataquian/news/internal/store"
	"github.com/riacataquian/news/web/cron/persistence"
)

const (
	domains Key = "domains"
	sources Key = "sources"
	query   Key = "query"

	defaultLang = "en"
)

var (
	client newsclient.HTTPClient

	timer        = clock.New()
	listEndpoint = list.ServiceEndpoint
)

// topQueried are hard-coded values that represents the top queried news given a domain.
// TODO: Replace me with actual values, retrieve from, ideally, a data repository.
var topQueried = []TopQueried{
	{
		Key:    domains,
		Values: []string{"techcrunch.com", "nytimes.com", "wsj.com"},
	},
	{
		Key:    sources,
		Values: []string{"bloomberg", "financial-times", "the-wall-street-journal"},
	},
	{
		Key:    query,
		Values: []string{"bitcoin", "ethereum", "blockchain"},
	},
}

// List fetches news as per topQueried values.
//
// It connects to https://newsapi.org to fetch the first 20 news
// that matches the values defined in topQueried, per key,
// and persist the results to the datastore.
// Finally, it returns the log containing the query parameters and the elapsed time
// performing the transactions.
//
// By default, language is set to "en".
// See https://newsapi.org/docs/endpoints/everything > Request Parameters
// on how to construct params.
//
// NOTE:
// Current newsapi plan fetch news anything not older than 7days from now.
// Future plans includes fetch all data which are 7 days old.
func List(ctx context.Context, repo store.Store, r *http.Request) (*Log, error) {
	// Requests to external services should timeout for 5 seconds.
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client = newsclient.NewFromContext(reqCtx, listEndpoint)
	started := timer.Now()

	var queried []TopQueried
	for _, top := range topQueried {
		switch domain := top.Key; domain {
		case domains:
			params := &list.Params{
				Language: defaultLang,
				Domains:  strings.Join(top.Values, ","),
			}
			_, err := fetchAndPersist(ctx, repo, client, params)
			if err != nil {
				return nil, err
			}
			queried = append(queried, top)
		case sources:
			params := &list.Params{
				Language: defaultLang,
				Sources:  strings.Join(top.Values, ","),
			}
			_, err := fetchAndPersist(ctx, repo, client, params)
			if err != nil {
				return nil, err
			}
			queried = append(queried, top)
		case query:
			// Surround phrases with quotes for exact match.
			var q []string
			for _, term := range top.Values {
				q = append(q, fmt.Sprintf("%q", term))
			}

			params := &list.Params{
				Language: defaultLang,
				Query:    strings.Join(q, "+"),
			}
			_, err := fetchAndPersist(ctx, repo, client, params)
			if err != nil {
				return nil, err
			}
			queried = append(queried, top)
		default:
			log.Printf("unknown domain: %v", domain)
		}
	}

	return &Log{
		Queried:     queried,
		ElapsedTime: timer.Since(started),
	}, nil
}

// fetchAndPersist connects to newsapi via a newsclient
// then persists the results to the supplied repo.
func fetchAndPersist(ctx context.Context, repo store.Store, client newsclient.HTTPClient, params newsclient.Params) (*news.Response, error) {
	authKey, err := auth.LookupAPIAuthKey()
	if err != nil {
		return nil, err
	}

	res, err := client.Get(authKey, params)
	if err != nil {
		return nil, err
	}

	if len(res.Articles) > 0 {
		err := persistence.Create(ctx, repo, timer, res.Articles)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return &news.Response{
		Status:       res.Status,
		TotalResults: 0,
	}, nil
}
