package cron

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/list"
	"github.com/riacataquian/news/internal/store"
)

var (
	originalClient     = client
	originalRepo       = repo
	originalTimer      = timer
	originalTopQueried = topQueried
)

func TestList(t *testing.T) {
	tests := []struct {
		desc          string
		isServerValid bool
		topQueried    []TopQueried
		wantLog       *Log
	}{
		{
			desc:          "returns the elapsed time after fetching topQueried items",
			isServerValid: true,
			wantLog: &Log{
				Queried:     topQueried,
				ElapsedTime: 123,
			},
		},
		{
			desc:          "does not query unknown domain",
			isServerValid: true,
			topQueried: []TopQueried{
				{
					Key:    Key("unknown-domain"),
					Values: []string{"some", "valid", "terms"},
				},
				{
					Key:    domains,
					Values: []string{"some", "valid", "terms"},
				},
			},
			wantLog: &Log{
				Queried: []TopQueried{
					{
						Key:    domains,
						Values: []string{"some", "valid", "terms"},
					},
				},
				ElapsedTime: 123,
			},
		},
	}

	for _, test := range tests {
		fakes, teardown := setup(t, config{isServerValid: true, isStoreValid: true})
		listEndpoint = newsclient.ServiceEndpoint{
			RequestURL: fakes.server.URL,
			DocsURL:    "http://fake-docs-url",
		}
		timer = fakeclock{nsec: 123}
		if len(test.topQueried) > 0 {
			topQueried = test.topQueried
		}
		defer teardown()

		r := httptest.NewRequest("GET", "/test", nil)
		got, err := List(context.Background(), r)
		if err != nil {
			t.Errorf("%s: List(_, _): want (_, nil), got (_, %v)", test.desc, err)
		}

		if test.wantLog != nil {
			if diff := pretty.Compare(got, test.wantLog); diff != "" {
				t.Errorf("%s: List(_, _) diff: (-got +want)\n%s", test.desc, diff)
			}
		}
	}
}

func TestFetchAndPersist(t *testing.T) {
	tests := []struct {
		desc          string
		params        *list.Params
		isStoreValid  bool
		isServerValid bool
		isClientValid bool
		withArticles  bool
		wantResponse  *news.Response
		wantRows      []store.Row
	}{
		{
			desc: "returns news.Response given list.Params",
			params: &list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			isStoreValid:  true,
			isServerValid: true,
			isClientValid: true,
			withArticles:  true,
			wantResponse:  fakeResponse,
		},
		{
			desc: "returns 0 results for news.Response given list.Params",
			params: &list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			isStoreValid:  true,
			isServerValid: true,
			isClientValid: true,
			wantResponse: &news.Response{
				Status:       "200",
				TotalResults: 0,
			},
		},
		{
			desc: "persists articles from news.Response given list.Params",
			params: &list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			isStoreValid:  true,
			isServerValid: true,
			isClientValid: true,
			withArticles:  true,
			wantRows: []store.Row{
				toStoreRow(
					123,
					"some-author-1",
					"some-title-1",
					"some-description-1",
					"some-URL-1",
					"some-image-url-1",
					time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				),
				toStoreRow(
					123,
					"some-author-2",
					"some-title-2",
					"some-description-2",
					"some-URL-2",
					"some-image-url-2",
					time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				),
				toStoreRow(
					123,
					"bloomberg",
					"Bloomberg",
				),
				toStoreRow(
					123,
					"financial-times",
					"Financial Times",
				),
			},
		},
	}

	for _, test := range tests {
		fakes, teardown := setup(t, config{
			isServerValid: test.isServerValid,
			isStoreValid:  test.isStoreValid,
			clockNanosec:  123,
		})
		client = &fakeclient{
			withArticles: test.withArticles,
			isValid:      test.isClientValid,
			serviceEndpoint: serviceEndpoint{
				RequestURL: fakes.server.URL,
			},
		}
		defer teardown()

		got, err := fetchAndPersist(context.Background(), client, test.params)
		if err != nil {
			t.Errorf("fetchAndPersist(_, _, %v): want (%v, nil), got (%v, %v)", test.params, test.wantResponse, got, err)
		}

		if test.wantResponse != nil {
			if diff := pretty.Compare(got, test.wantResponse); diff != "" {
				t.Errorf("%s: fetchAndPersist(_, _, %v) diff: (-got +want)\n%s", test.desc, test.params, diff)
			}
		}

		if len(test.wantRows) > 0 {
			if diff := pretty.Compare(fakes.store.rows, test.wantRows); diff != "" {
				t.Errorf("%s: fetchAndPersist(_, _, %v) diff: (-got +want)\n%s", test.desc, test.params, diff)
			}
		}
	}
}

func TestFetchAndPersistErrors(t *testing.T) {
	tests := []struct {
		desc          string
		params        *list.Params
		isServerValid bool
		isStoreValid  bool
		isClientValid bool
		withArticles  bool
	}{
		{
			desc:          "returns an error when client errored",
			isServerValid: true,
			isStoreValid:  true,
			isClientValid: false,
		},
		{
			desc:          "returns an error when store errored",
			isServerValid: true,
			isStoreValid:  false,
			isClientValid: true,
			withArticles:  true,
		},
	}

	for _, test := range tests {
		fakes, teardown := setup(t, config{
			isServerValid: test.isServerValid,
			isStoreValid:  test.isStoreValid,
		})
		client = &fakeclient{
			isValid:      test.isClientValid,
			withArticles: test.withArticles,
			serviceEndpoint: serviceEndpoint{
				RequestURL: fakes.server.URL,
			},
		}
		defer teardown()

		params := &list.Params{
			Language: defaultLang,
			Domains:  "some-domain-1,some-domain-2",
		}
		got, err := fetchAndPersist(context.Background(), client, params)
		if err == nil {
			t.Errorf("%s: fetchAndPersist(_, _, %v): want (_, error), got (%v, %v)", test.desc, test.params, got, err)
		}
	}
}
