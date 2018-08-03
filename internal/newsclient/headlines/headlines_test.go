package headlines

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
)

// FakeClient mocks a Client interface.
type FakeClient struct {
	newsclient.ServiceEndpoint
	RequestOrigin *http.Request
	IsValid       bool
}

func (f FakeClient) Get(_ context.Context, _ *http.Request, p Params) (*news.Response, error) {
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

func (f FakeClient) DispatchRequest(_ context.Context, r *http.Request) (*news.Response, error) {
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

func setupAPIKey(t *testing.T) {
	t.Helper()

	err := os.Setenv("API_KEY", "this is a test api key")
	if err != nil {
		t.Fatal(err)
	}
}

func setupStubServer(t *testing.T, isValid bool) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isValid {
			errString := `{"status": "internal server error", "code": "500", "message": "some error"}`
			http.Error(w, errString, http.StatusNotFound)
			return
		}

		resp := &news.Response{
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
		}
		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("error marshalling response: %v", err)
		}
		w.Write(b)
	}))
}

func TestGet(t *testing.T) {
	setupAPIKey(t)

	want := &news.Response{
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
	}

	server := setupStubServer(t, true)
	defer server.Close()

	client := FakeClient{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: server.URL,
		},
		IsValid: true,
	}

	ctx := context.Background()
	r := httptest.NewRequest("GET", server.URL, nil)
	got, err := client.Get(ctx, r, Params{Country: "us"})
	if err != nil {
		t.Errorf("Get: want (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a news.Response and nil error"
		t.Errorf("%s: Get diff: (-got +want)\n%s", desc, diff)
	}
}

func TestGetErrors(t *testing.T) {
	err := os.Setenv("API_KEY", "this is a test api key")
	if err != nil {
		t.Logf("Get: setting up an API_KEY: %v", err)
	}

	tests := []struct {
		desc          string
		isServerValid bool
		isClientValid bool
		params        Params
	}{
		{
			desc:          "returns an error when server errored",
			isServerValid: false,
			isClientValid: true,
			params:        Params{Country: "us"},
		},
		{
			desc:          "returns an error when client errored",
			isServerValid: true,
			isClientValid: false,
			params:        Params{Country: "us"},
		},
		{
			desc:          "returns an error when params errored",
			isServerValid: true,
			isClientValid: true,
			params:        Params{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			server := setupStubServer(t, false)
			defer server.Close()

			r := httptest.NewRequest("GET", server.URL, nil)
			ctx := context.Background()
			client := FakeClient{
				ServiceEndpoint: newsclient.ServiceEndpoint{
					URL: server.URL,
				},
			}
			got, err := client.Get(ctx, r, test.params)
			if err == nil {
				t.Errorf("%s: Get(_, _, %v) want (nil, error), got (%v, %v)", test.desc, test.params, got, err)
			}
		})
	}
}

func TestDispatchRequest(t *testing.T) {
	want := &news.Response{
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
	}

	server := setupStubServer(t, true)
	defer server.Close()

	r, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("DispatchRequest(_, _): error creating a new request: %v", err)
	}

	got, err := Client{}.DispatchRequest(context.Background(), r)
	if err != nil {
		t.Errorf("DispatchRequest(_, _): want (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a news.Response and nil error"
		t.Errorf("%s: DispatchRequest(_, _) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestDispatchRequestErrors(t *testing.T) {
	want := &news.ErrorResponse{
		Status:  "internal server error",
		Code:    "500",
		Message: "some error",
	}

	server := setupStubServer(t, false)
	defer server.Close()

	r, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("DispatchRequest(_, _): error creating a new request: %v", err)
	}

	got, err := Client{}.DispatchRequest(context.Background(), r)
	if err == nil {
		t.Errorf("DispatchRequest(_, _): want (nil, error), got (%v, %v)", got, err)
	}

	if diff := pretty.Compare(err, want); diff != "" {
		desc := "returns a news.ErrorResponse when error is encountered"
		t.Errorf("%s: DispatchRequest(_, _) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		desc string
		in   Params
		want string
	}{
		{
			desc: "returns the encoded params",
			in:   Params{Country: "us"},
			want: "country=us",
		},
		{
			desc: "returns the correct query params",
			in:   Params{Country: "us", Query: "bitcoin", PageSize: 50, Page: 2},
			want: "country=us&page=2&pageSize=50&q=bitcoin",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := test.in.Encode()
			if got != test.want {
				t.Errorf("Encode: want (%v, nil), got (%v, %v)", test.want, got, err)
			}
		})
	}
}

func TestEncodeErrors(t *testing.T) {
	tests := []struct {
		desc string
		in   Params
	}{
		{
			desc: "country can't be mixed with sources param",
			in:   Params{Country: "us", Sources: "the-times-of-india"},
		},
		{
			desc: "category can't be mixed with sources param",
			in:   Params{Category: "technology", Sources: "the-times-of-india"}},
		{
			desc: "pageSize exceeded the maxPageSize",
			in:   Params{PageSize: 500, Category: "technology", Sources: "the-times-of-india"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := test.in.Encode()
			if err == nil {
				t.Errorf("%s: (%v).Encode(): want (nil, error), got (%v, %v)", test.desc, test.in, got, err)
			}
		})
	}
}
