package persistence

// This file contains fake implementations used for testing.

import (
	"errors"
	"time"

	"github.com/riacataquian/news/internal/store"
)

type fakestore struct {
	isValid bool
	// rows are the supposedly inserted rows.
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
