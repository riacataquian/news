package list

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
)

func setupAPIKey(t *testing.T) {
	t.Helper()

	err := os.Setenv("API_KEY", "this is a test api key")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewClient(t *testing.T) {
	got := NewClient()
	want := &Client{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: Endpoint,
		},
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a new list newsclient"
		t.Errorf("%s: NewClient(): Diff (-got +want)\n%s", desc, diff)
	}
}

func TestGet(t *testing.T) {
	setupAPIKey(t)

	want := &news.Response{
		Status:       "200",
		TotalResults: 2,
		Articles: []*news.News{
			{
				Source: &news.Source{
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
				Source: &news.Source{
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

	client := fakeclient{isValid: true}

	ctx := context.Background()
	got, err := client.Get(ctx, Params{SortBy: Relevancy, Language: "en"})
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
			params:        Params{SortBy: Relevancy, Language: "en"},
		},
		{
			desc:          "returns an error when client errored",
			isServerValid: true,
			isClientValid: false,
			params:        Params{SortBy: Relevancy, Language: "en"},
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

			ctx := context.Background()
			client := fakeclient{}
			got, err := client.Get(ctx, test.params)
			if err == nil {
				t.Errorf("%s: Get(_, %v) want (nil, error), got (%v, %v)", test.desc, test.params, got, err)
			}
		})
	}
}

func TestDispatchRequest(t *testing.T) {
	want := &news.Response{
		Status:       "200",
		TotalResults: 2,
		Articles: []*news.News{
			{
				Source: &news.Source{
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
				Source: &news.Source{
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
		t.Fatalf("DispatchRequest(_): error creating a new request: %v", err)
	}

	client := &Client{}
	got, err := client.DispatchRequest(context.Background(), r)
	if err != nil {
		t.Errorf("DispatchRequest(_): want (%v, nil), got (%v, %v)", want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a news.Response and nil error"
		t.Errorf("%s: DispatchRequest(_) diff: (-got +want)\n%s", desc, diff)
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
		t.Fatalf("DispatchRequest(_): error creating a new request: %v", err)
	}

	client := &Client{}
	got, err := client.DispatchRequest(context.Background(), r)
	if err == nil {
		t.Errorf("DispatchRequest(_): want (nil, error), got (%v, %v)", got, err)
	}

	if diff := pretty.Compare(err, want); diff != "" {
		desc := "returns a news.ErrorResponse when error is encountered"
		t.Errorf("%s: DispatchRequest(_) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestEncode(t *testing.T) {
	in := Params{
		Query:    "some-query",
		Sources:  "some-source1,some-source2",
		Domains:  "some-domain1,some-domain2",
		SortBy:   Popularity,
		Language: "en",
	}
	want := "domains=some-domain1%2Csome-domain2&language=en&q=some-query&sortBy=popularity&sources=some-source1%2Csome-source2"
	got, err := in.Encode()
	if got != want {
		t.Errorf("Encode: want (%v, nil), got (%v, %v)", want, got, err)
	}
}

func TestEncodeErrors(t *testing.T) {
	in := Params{PageSize: 500, Language: "en"}
	got, err := in.Encode()
	if err == nil {
		desc := "pageSize exceeded the maxPageSize"
		t.Errorf("%s: (%v).Encode(): want (nil, error), got (%v, %v)", desc, in, got, err)
	}
}
