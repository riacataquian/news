package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/store"
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
	isError bool
	serviceEndpoint
}

func (f *fakeclient) Get(ctx context.Context, authKey string, p newsclient.Params) (*news.Response, error) {
	if f.isError {
		return nil, errors.New("some error")
	}

	_, err := p.Encode()
	if err != nil {
		return nil, err
	}

	return fakeResponse, nil
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

func setupStubServer(t *testing.T, isError bool) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isError {
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
	isServerError bool
	isClientError bool
}

type fakestore struct{}

func (f *fakestore) Create(_ string, _ []string, _ ...store.Row) error {
	return nil
}

// fakes encapsulates a test's fake structures.
type fakes struct {
	server *httptest.Server
	store  *fakestore
	client newsclient.HTTPClient
}

// setup performs the necessary monkey-patching per test suite
// then return a teardown function to return back original values of package wide vars.
func setup(t *testing.T, conf config) (*fakes, teardown) {
	t.Helper()

	os.Setenv("API_KEY", "test-api-key")

	fakeserver := setupStubServer(t, conf.isServerError)
	fakeclient := &fakeclient{
		isError: conf.isClientError,
		serviceEndpoint: serviceEndpoint{
			RequestURL: fakeserver.URL,
		},
	}
	fakestore := &fakestore{}
	client = fakeclient

	fakes := fakes{
		server: fakeserver,
		store:  fakestore,
		client: fakeclient,
	}

	teardown := func() {
		os.Clearenv()
		fakeserver.Close()

		client = originalClient

		headlinesEndpoint = originalHeadlinesEndpoint
		listEndpoint = originalListEndpoint
	}

	return &fakes, teardown
}
