package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"

	"github.com/kylelemons/godebug/pretty"
)

var originalClient = client

// FakeParams ...
type FakeParams string

// Encode ...
func (fp FakeParams) Encode() (string, error) {
	return "sources=bloomberg,financial-times", nil
}

// FakeClient mocks a newsclient.Client interface.
type FakeClient struct {
	newsclient.ServiceEndpoint
	ContextOrigin context.Context
	RequestOrigin *http.Request
	IsValid       bool
}

func (f FakeClient) GetContextOrigin() context.Context {
	return f.ContextOrigin
}

func (f FakeClient) GetRequestOrigin() *http.Request {
	return f.RequestOrigin
}

func (f FakeClient) GetServiceEndpoint() newsclient.ServiceEndpoint {
	return f.ServiceEndpoint
}

func (f FakeClient) GetTopHeadlines(p newsclient.Params) (*news.Response, error) {
	if f.IsValid {
		return nil, errors.New("some error")
	} else {
		return &news.Response{
			Status:       "200",
			TotalResults: 2,
			Articles: []news.News{
				{
					Source: news.Source{
						ID:   "bloomberg",
						Name: "Bloomberg",
					},
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
				{
					Source: news.Source{
						ID:   "financial-times",
						Name: "Financial Times",
					},
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil
	}
}

func (f FakeClient) DispatchRequest(r *http.Request) (*news.Response, error) {
	if f.IsValid {
		return &news.Response{
			Status:       "200",
			TotalResults: 2,
			Articles: []news.News{
				{
					Source: news.Source{
						ID:   "bloomberg",
						Name: "Bloomberg",
					},
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil
	}

	return nil, errors.New("failed request")
}

func setupFakeClient(t *testing.T, queryParams string, isValid bool) FakeClient {
	t.Helper()

	ctx := context.Background()
	r := httptest.NewRequest("GET", "/test?"+queryParams, nil)
	return FakeClient{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: "test-url",
		},
		ContextOrigin: ctx,
		RequestOrigin: r,
		IsValid:       isValid,
	}
}

func teardown(t *testing.T) {
	t.Helper()
	client = originalClient
}

func TestFetchNews(t *testing.T) {
	qp := "sources=bloomberg,financial-times"
	client = setupFakeClient(t, qp, false)
	defer teardown(t)

	want := &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: "/test?sources=bloomberg,financial-times",
		Data: &news.Response{
			Status:       "200",
			TotalResults: 2,
			Articles: []news.News{
				{
					Source: news.Source{
						ID:   "bloomberg",
						Name: "Bloomberg",
					},
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
				{
					Source: news.Source{
						ID:   "financial-times",
						Name: "Financial Times",
					},
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	params := FakeParams(qp)
	got, err := fetchNews(client, params)
	if err != nil {
		t.Errorf("fetchNews: expecting (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a SuccessResponse with correct Code and RequestURL"
		t.Errorf("%s: fetchNews diff: (-got +want)\n%s", desc, diff)
	}
}

func TestFetchNewsErrors(t *testing.T) {
	qp := "sources=bloomberg,financial-times"
	client = setupFakeClient(t, qp, true)
	defer teardown(t)

	params := FakeParams(qp)
	got, err := fetchNews(client, params)
	if err == nil {
		desc := "returns nil SuccessResponse when an error is encountered"
		t.Errorf("%s: fetchNews expecting (nil, error), got (%v, %v)", desc, got, err)
	}
}
