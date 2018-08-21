package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
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

// fakeclient mocks a newsclient.Client interface.
type fakeclient struct {
	isValid bool
	serviceEndpoint
}

func (f *fakeclient) Get(authKey string, p newsclient.Params) (*news.Response, error) {
	_, err := p.Encode()
	if err != nil {
		return nil, err
	}

	if f.isValid {
		return fakeResponse, nil
	}
	return nil, errors.New("some error")
}

// fakeParams mocks a newsclient.Params interface.
type fakeParams struct {
	lang    string
	wantErr bool
}

func (fp fakeParams) Encode() (string, error) {
	if fp.wantErr {
		return "", errors.New("error encoding params")
	}

	return "lang=" + fp.lang, nil
}

func (fp fakeParams) Read(_ []byte) (n int, err error) {
	return len(fp.lang), nil
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

// serviceEndpoint embeds a RequestURL.
type serviceEndpoint struct {
	RequestURL string
}

type teardown func()

// config encapsulates a test's setup configuration.
type config struct {
	isServerValid bool
	isClientValid bool
}

// fakes encapsulates a test's fake structures.
type fakes struct {
	server *httptest.Server
	client newsclient.HTTPClient
}

// setup performs the necessary monkey-patching per test suite
// then return a teardown function to return back original values of package wide vars.
func setup(t *testing.T, conf config) (*fakes, teardown) {
	t.Helper()

	os.Setenv("API_KEY", "test-api-key")

	fakeserver := setupStubServer(t, conf.isServerValid)
	fakeclient := &fakeclient{
		isValid: conf.isClientValid,
		serviceEndpoint: serviceEndpoint{
			RequestURL: fakeserver.URL,
		},
	}

	client = fakeclient

	fakes := fakes{
		server: fakeserver,
		client: fakeclient,
	}

	teardown := func() {
		os.Setenv("API_KEY", "")
		fakeserver.Close()

		client = originalClient

		headlinesEndpoint = originalHeadlinesEndpoint
		listEndpoint = originalListEndpoint
	}

	return &fakes, teardown
}
