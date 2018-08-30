package handler

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/internal/newsclient"
)

var (
	originalClient            = client
	originalDefaultDuration   = defaultDuration
	originalHeadlinesEndpoint = headlinesEndpoint
	originalListEndpoint      = listEndpoint
)

func TestList(t *testing.T) {
	fakes, teardown := setup(t, config{})
	listEndpoint = newsclient.ServiceEndpoint{
		RequestURL: fakes.server.URL,
		DocsURL:    "http://fake-docs-url",
	}
	defer teardown()

	req, err := http.NewRequest(http.MethodGet, fakes.server.URL, nil)
	req.Form = url.Values{"query": {"bitcoin"}, "page": {"2"}}
	if err != nil {
		t.Fatalf("List(_, _): got error: %v, want nil error", err)
	}

	want := &SuccessResponse{
		Code:       http.StatusOK,
		Count:      len(fakeResponse.Articles),
		Page:       2,
		TotalCount: fakeResponse.TotalResults,
		Data:       fakeResponse.Articles,
	}
	desc := "returns the list of news given query parameter"
	got, err := List(context.Background(), fakes.store, req)
	if err != nil {
		t.Fatalf("%s: List(_, _, _): want(%v, nil), got (%v, %v)", desc, want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("%s: List(_, _, _) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestListErrors(t *testing.T) {
	tests := []struct {
		desc          string
		isServerError bool
		params        url.Values
	}{
		{
			desc:          "returns an error when server errored",
			isServerError: true,
			params:        url.Values{"query": {"valid-query"}},
		},
		{
			desc:   "returns an error when encoding params errored",
			params: url.Values{"unrecognized-key": {"unrecognized-value"}},
		},
	}

	for _, test := range tests {
		fakes, teardown := setup(t, config{isServerError: test.isServerError})
		listEndpoint = newsclient.ServiceEndpoint{
			RequestURL: fakes.server.URL,
			DocsURL:    "http://fake-docs-url",
		}
		defer teardown()

		req, err := http.NewRequest(http.MethodGet, fakes.server.URL, nil)
		req.Form = test.params
		if err != nil {
			t.Fatalf("List(_, _, _): got error: %v, want nil error", err)
		}

		if got, err := List(context.Background(), fakes.store, req); err == nil {
			t.Errorf("%s: List(_, _, _), expecting (nil, error), got (%v, %v)", test.desc, got, err)
		}
	}
}

func TestTopHeadlines(t *testing.T) {
	fakes, teardown := setup(t, config{})
	headlinesEndpoint = newsclient.ServiceEndpoint{
		RequestURL: fakes.server.URL,
		DocsURL:    "http://fake-docs-url",
	}
	defer teardown()

	req, err := http.NewRequest(http.MethodGet, fakes.server.URL, nil)
	req.Form = url.Values{"query": {"bitcoin"}}
	if err != nil {
		t.Fatalf("TopHeadlines(_, _, _): got error: %v, want nil error", err)
	}

	want := &SuccessResponse{
		Code:       http.StatusOK,
		Count:      len(fakeResponse.Articles),
		Page:       1,
		TotalCount: fakeResponse.TotalResults,
		Data:       fakeResponse.Articles,
	}
	desc := "returns the top headlines news given query parameter"
	got, err := TopHeadlines(context.Background(), fakes.store, req)
	if err != nil {
		t.Fatalf("%s: TopHeadlines(_, _, _): want(%v, nil), got (%v, %v)", desc, want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("%s: TopHeadlines(_, _, _) diff: (-got +want)\n%s", desc, diff)
	}
}

func TestTopHeadlinesErrors(t *testing.T) {
	tests := []struct {
		desc          string
		isServerError bool
		params        url.Values
	}{
		{
			desc:          "returns an error when server errored",
			isServerError: true,
			params:        url.Values{"query": {"bitcoin"}},
		},
		{
			desc:   "returns an error when encoding params errored",
			params: url.Values{"unrecognized-key": {"unrecognized-value"}},
		},
	}

	for _, test := range tests {
		fakes, teardown := setup(t, config{isServerError: test.isServerError})
		headlinesEndpoint = newsclient.ServiceEndpoint{
			RequestURL: fakes.server.URL,
			DocsURL:    "http://fake-docs-url",
		}
		defer teardown()

		req, err := http.NewRequest(http.MethodGet, fakes.server.URL, nil)
		req.Form = test.params
		if err != nil {
			t.Fatalf("TopHeadlines(_, _, _): got error: %v, want nil error", err)
		}

		if got, err := TopHeadlines(context.Background(), fakes.store, req); err == nil {
			t.Errorf("%s: TopHeadlines(_, _, _), expecting (nil, error), got (%v, %v)", test.desc, got, err)
		}
	}
}

func TestFetch(t *testing.T) {
	_, teardown := setup(t, config{})
	defer teardown()

	desc := "returns a SuccessResponse with correct Code and RequestURL"
	params := fakeParams{lang: "en"}
	want := fakeResponse
	got, err := fetch(params)
	if err != nil {
		t.Fatalf("%s: fetch(%v): expecting (%v, nil), got (%v, %v)", desc, params, want, got, err)
	}

	if diff := pretty.Compare(got, want); diff != "" {
		t.Errorf("%s: fetch(%v), diff: (-got +want)\n%s", desc, params, diff)
	}
}

func TestFetchErrors(t *testing.T) {
	tests := []struct {
		desc          string
		params        *fakeParams
		isClientError bool
	}{
		{
			desc:          "returns an error when an error in client is encountered",
			params:        &fakeParams{},
			isClientError: true,
		},
		{
			desc:   "returns an error when an error in params is encountered",
			params: &fakeParams{wantErr: true},
		},
	}

	for _, test := range tests {
		_, teardown := setup(t, config{isClientError: test.isClientError})
		defer teardown()

		if got, err := fetch(test.params); err == nil {
			t.Errorf("%s: fetch(%v), expecting (nil, error), got (%v, %v)", test.desc, test.params, got, err)
		}
	}
}
