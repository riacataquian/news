package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"

	"github.com/kylelemons/godebug/pretty"
)

var originalClient = client

type config struct {
	ctx         context.Context
	req         *http.Request
	queryParams string
	isValid     bool
}

// FakeParams mocks a newsclient.Params interface.
type FakeParams string

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

func (f FakeClient) TopHeadlines(_ context.Context, _ *http.Request, p newsclient.Params) (*news.Response, error) {
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
	return nil, errors.New("some error")
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

func setupFakeClient(t *testing.T, conf config) FakeClient {
	t.Helper()

	return FakeClient{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: "test-url",
		},
		ContextOrigin: conf.ctx,
		RequestOrigin: conf.req,
		IsValid:       conf.isValid,
	}
}

func teardown(t *testing.T) {
	t.Helper()
	client = originalClient
}

func TestFetchNews(t *testing.T) {
	q := "sources=bloomberg,financial-times"
	conf := config{
		ctx:         context.Background(),
		req:         httptest.NewRequest("GET", fmt.Sprintf("/test?%s", q), nil),
		queryParams: q,
		isValid:     true,
	}
	client = setupFakeClient(t, conf)
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

	got, err := fetchNews(conf.ctx, conf.req, client, FakeParams(conf.queryParams))
	if err != nil {
		t.Errorf("fetchNews: expecting (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a SuccessResponse with correct Code and RequestURL"
		t.Errorf("%s: fetchNews diff: (-got +want)\n%s", desc, diff)
	}
}

func TestFetchNewsErrors(t *testing.T) {
	q := "sources=bloomberg,financial-times"
	conf := config{
		ctx:         context.Background(),
		req:         httptest.NewRequest("GET", fmt.Sprintf("/test?%s", q), nil),
		queryParams: q,
		isValid:     false,
	}
	client = setupFakeClient(t, conf)
	defer teardown(t)

	got, err := fetchNews(conf.ctx, conf.req, client, FakeParams(conf.queryParams))
	if err == nil {
		desc := "returns nil SuccessResponse when an error is encountered"
		t.Errorf("%s: fetchNews expecting (nil, error), got (%v, %v)", desc, got, err)
	}
}
