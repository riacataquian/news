package cron

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/list"
	"github.com/riacataquian/news/internal/store"
)

type fakeclient struct {
	withArticles bool
	isValid      bool
	params       []list.Params
}

func (f *fakeclient) Get(_ context.Context, p newsclient.Params) (*news.Response, error) {
	f.params = append(f.params, p.(list.Params))

	if !f.isValid {
		return nil, errors.New("some error")
	}

	if f.withArticles {
		return &news.Response{
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
		}, nil
	}

	return &news.Response{Status: "200", TotalResults: 0}, nil
}

func (f *fakeclient) DispatchRequest(_ context.Context, r *http.Request) (*news.Response, error) {
	if !f.isValid {
		return nil, errors.New("some error")
	}

	if f.withArticles {
		return &news.Response{
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
					URL:         "some-URL",
					ImageURL:    "some-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil
	}

	return &news.Response{Status: "200", TotalResults: 0}, nil
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
