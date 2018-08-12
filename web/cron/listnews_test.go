package cron

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient/list"
	"github.com/riacataquian/news/internal/store"
)

var (
	originalClient     = listclient
	originalRepo       = repo
	originalTimer      = timer
	originalTopQueried = topQueried
)

func setupAPIKey(t *testing.T) {
	t.Helper()

	err := os.Setenv("API_KEY", "this is a test api key")
	if err != nil {
		t.Fatal(err)
	}
}

func teardown() {
	listclient = originalClient
	repo = originalRepo
	timer = originalTimer
	topQueried = originalTopQueried
}

func TestList(t *testing.T) {
	setupAPIKey(t)

	tests := []struct {
		desc       string
		client     *fakeclient
		topQueried []TopQueried
		wantLog    *Log
		wantParams []list.Params
		wantRows   []store.Row
	}{
		{
			desc:   "returns the elapsed time after fetching topQueried items",
			client: &fakeclient{isValid: true},
			wantLog: &Log{
				Queried:     topQueried,
				ElapsedTime: 123,
			},
		},
		{
			desc:   "does not query unknown domain",
			client: &fakeclient{isValid: true},
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
		{
			desc:   "contructs list.Params based on topQueried values for querying news",
			client: &fakeclient{isValid: true},
			topQueried: []TopQueried{
				{
					Key:    domains,
					Values: []string{"test-domain-1", "test-domain-2"},
				},
				{
					Key:    sources,
					Values: []string{"test-source-1", "test-source-2"},
				},
				{
					Key:    query,
					Values: []string{"test-query-1", "test-query-2"},
				},
			},
			wantParams: []list.Params{
				{
					Language: defaultLang,
					Domains:  "test-domain-1,test-domain-2",
				},
				{
					Language: defaultLang,
					Sources:  "test-source-1,test-source-2",
				},
				{
					Language: defaultLang,
					Query:    `"test-query-1"+"test-query-2"`,
				},
			},
		},
	}

	for _, test := range tests {
		if len(test.topQueried) > 0 {
			topQueried = test.topQueried
		}

		listclient = test.client
		repo = func() store.Store {
			return &fakestore{isValid: true}
		}
		timer = fakeclock{nsec: 123}
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

		if len(test.wantParams) > 0 {
			if diff := pretty.Compare(test.client.params, test.wantParams); diff != "" {
				t.Errorf("%s: List(_, _) diff: (-got +want)\n%s", test.desc, diff)
			}
		}
	}
}

func TestFetchAndPersist(t *testing.T) {
	setupAPIKey(t)

	tests := []struct {
		desc         string
		params       list.Params
		store        *fakestore
		client       *fakeclient
		wantResponse *news.Response
		wantRows     []store.Row
	}{
		{
			desc: "returns news.Response given list.Params",
			params: list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			store: &fakestore{isValid: true},
			client: &fakeclient{
				isValid:      true,
				withArticles: true,
			},
			wantResponse: &news.Response{
				Status:       "200",
				TotalResults: 1,
				Articles: []*news.News{
					{
						Source: &news.Source{
							ID:   "some-source-id",
							Name: "some-source-name",
						},
						Author:      "some-author",
						Title:       "some-title",
						Description: "some-description",
						URL:         "http://test-url",
						ImageURL:    "http://test-image-url",
						PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 123, time.UTC),
					},
				},
			},
		},
		{
			desc: "returns 0 results for news.Response given list.Params",
			params: list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			store:  &fakestore{isValid: true},
			client: &fakeclient{isValid: true},
			wantResponse: &news.Response{
				Status:       "200",
				TotalResults: 0,
			},
		},
		{
			desc: "persists articles from news.Response given list.Params",
			params: list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			store: &fakestore{isValid: true},
			client: &fakeclient{
				isValid:      true,
				withArticles: true,
			},
			wantRows: []store.Row{
				toStoreRow(
					123,
					"some-author",
					"some-title",
					"some-description",
					"http://test-url",
					"http://test-image-url",
					time.Date(2016, time.August, 15, 0, 0, 0, 123, time.UTC),
				),
				toStoreRow(
					123,
					"some-source-id",
					"some-source-name",
				),
			},
		},
	}

	for _, test := range tests {
		listclient = test.client
		timer = &fakeclock{nsec: 123}
		repo = func() store.Store {
			return test.store
		}
		defer teardown()

		r := httptest.NewRequest("GET", "/test", nil)
		got, err := fetchAndPersist(context.Background(), r, test.params)
		if err != nil {
			t.Errorf("fetchAndPersist(_, _, %v): want (%v, nil), got (%v, %v)", test.params, test.wantResponse, got, err)
		}

		if test.wantResponse != nil {
			if diff := pretty.Compare(got, test.wantResponse); diff != "" {
				t.Errorf("%s: fetchAndPersist(_, _, %v) diff: (-got +want)\n%s", test.desc, test.params, diff)
			}
		}

		if len(test.wantRows) > 0 {
			if diff := pretty.Compare(test.store.rows, test.wantRows); diff != "" {
				t.Errorf("%s: fetchAndPersist(_, _, %v) diff: (-got +want)\n%s", test.desc, test.params, diff)
			}
		}
	}
}

func TestFetchAndPersistErrors(t *testing.T) {
	setupAPIKey(t)

	tests := []struct {
		desc   string
		params list.Params
		client *fakeclient
		store  *fakestore
	}{
		{
			desc: "returns an error when client errored",
			params: list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			client: &fakeclient{isValid: false},
			store:  &fakestore{isValid: true},
		},
		{
			desc: "returns an error when store errored",
			params: list.Params{
				Language: defaultLang,
				Domains:  "some-domain-1,some-domain-2",
			},
			client: &fakeclient{isValid: true, withArticles: true},
			store:  &fakestore{isValid: false},
		},
	}

	for _, test := range tests {
		listclient = test.client
		repo = func() store.Store {
			return test.store
		}
		defer teardown()

		r := httptest.NewRequest("GET", "/test", nil)
		got, err := fetchAndPersist(context.Background(), r, test.params)
		if err == nil {
			t.Errorf("%s: fetchAndPersist(_, _, %v): want (_, error), got (%v, %v)", test.desc, test.params, got, err)
		}
	}
}
