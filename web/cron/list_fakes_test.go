package cron

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
			Author:      "some-author-1",
			Title:       "some-title-1",
			Description: "some-description-1",
			URL:         "some-URL-1",
			ImageURL:    "some-image-url-1",
			PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Source: &news.Source{
				ID:   "financial-times",
				Name: "Financial Times",
			},
			Author:      "some-author-2",
			Title:       "some-title-2",
			Description: "some-description-2",
			URL:         "some-URL-2",
			ImageURL:    "some-image-url-2",
			PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
		},
	},
}

// fakeclient mocks a newsclient.Client interface.
type fakeclient struct {
	isValid      bool
	withArticles bool
	serviceEndpoint
}

func (f *fakeclient) Get(authKey string, p newsclient.Params) (*news.Response, error) {
	_, err := p.Encode()
	if err != nil {
		return nil, err
	}
	if f.withArticles {
		return fakeResponse, nil
	}

	if f.isValid {
		return &news.Response{Status: "200", TotalResults: 0}, nil
	}

	return nil, errors.New("some error")
}

type fakestore struct {
	isValid bool
	// rows are the supposedly inserted rows.
	// rows is set after calling fakestore's Create method.
	rows []store.Row
}

func (f *fakestore) Create(table string, cols []string, rows ...store.Row) error {
	if f.isValid {
		f.rows = append(f.rows, rows...)
		return nil
	}

	return errors.New("some store error")
}

type fakeclock struct {
	nsec int
}

func (c fakeclock) Now() time.Time {
	return time.Date(2016, time.August, 15, 0, 0, 0, c.nsec, time.UTC)
}

func (c fakeclock) Since(_ time.Time) time.Duration {
	return 123
}

func toStoreRow(args ...interface{}) store.Row {
	var r store.Row
	for _, arg := range args {
		r = append(r, arg)
	}
	return r
}

func setupStubServer(t *testing.T, isServerValid bool) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isServerValid {
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

// config encapsulates a test's setup configuration.
type config struct {
	isServerValid bool
	isStoreValid  bool
	clockNanosec  int
}

// fakes encapsulates a test's fake structures.
type fakes struct {
	server *httptest.Server
	store  *fakestore
	clock  *fakeclock
}

type teardown func()

// setup performs the necessary monkey-patching per test suite
// then return a teardown function to return back original values of package wide vars.
func setup(t *testing.T, conf config) (*fakes, teardown) {
	t.Helper()

	os.Setenv("API_KEY", "test-api-key")

	fakeserver := setupStubServer(t, conf.isServerValid)
	fakestore := &fakestore{isValid: conf.isStoreValid}
	fakeclock := &fakeclock{nsec: conf.clockNanosec}

	timer = fakeclock
	repo = func() store.Store {
		return fakestore
	}

	fakes := fakes{
		server: fakeserver,
		store:  fakestore,
		clock:  fakeclock,
	}

	teardown := func() {
		os.Setenv("API_KEY", "")
		fakeserver.Close()

		client = originalClient
		repo = originalRepo
		timer = originalTimer
		topQueried = originalTopQueried
	}

	return &fakes, teardown
}
