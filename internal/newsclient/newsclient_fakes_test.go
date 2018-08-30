package newsclient

// This file contains fake definitions of client, response and server.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/riacataquian/news/api/news"
)

var fakeResponse = &news.Response{
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

// setupFakeClient returns a fakeClient.
// Supply `url` with a stub server's URL.
func setupFakeClient(url string) *Client {
	return &Client{
		ServiceEndpoint: ServiceEndpoint{
			RequestURL: url,
			DocsURL:    "some-docs-url",
		},
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

		b, err := json.Marshal(fakeResponse)
		if err != nil {
			t.Fatalf("error marshalling response: %v", err)
		}
		w.Write(b)
	}))
}
