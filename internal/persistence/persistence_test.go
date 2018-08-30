package persistence

import (
	"os"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/store"
)

func setupAPIKey(t *testing.T) {
	t.Helper()

	err := os.Setenv("API_KEY", "this is a test api key")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPersist(t *testing.T) {
	setupAPIKey(t)

	tests := []struct {
		desc  string
		clock *fakeclock
		repo  *fakestore
		want  []store.Row
		in    *News
	}{
		{
			desc:  "persists news rows to repo",
			clock: &fakeclock{nsec: 123},
			repo:  &fakestore{isValid: true},
			want: []store.Row{
				toStoreRow(
					123,
					"some-author",
					"some-title",
					"some-description",
					"http://test-url",
					"http://test-image-url",
					time.Date(2016, time.August, 15, 0, 0, 0, 123, time.UTC),
				),
			},
			in: &News{
				News: &news.News{
					Author:      "some-author",
					Title:       "some-title",
					Description: "some-description",
					URL:         "http://test-url",
					ImageURL:    "http://test-image-url",
					PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 123, time.UTC),
				},
			},
		},
		{
			desc:  "persists news rows and its sources to repo",
			clock: &fakeclock{nsec: 123},
			repo:  &fakestore{isValid: true},
			want: []store.Row{
				toStoreRow(
					123,
					"some-author",
					"some-title",
					"some-description",
					"http://test-url",
					"http://test-image-url",
					time.Date(2016, time.August, 15, 0, 0, 0, 123, time.UTC),
				),
				toStoreRow(
					123,
					"some-source-id",
					"some-source-name",
				),
			},
			in: &News{
				News: &news.News{
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
		},
	}

	for _, test := range tests {
		err := test.in.Create(test.repo, test.clock)
		if err != nil {
			t.Errorf("Create(_, %v): want nil, got %v", test.in, err)
		}

		if diff := pretty.Compare(test.repo.rows, test.want); diff != "" {
			t.Errorf("%s: Create(_, %v): Diff (-got +want)\n%s", test.desc, test.in, diff)
		}
	}
}

func TestCreateError(t *testing.T) {
	repo := &fakestore{isValid: false}
	clock := fakeclock{}
	in := &News{
		News: &news.News{
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
	}

	err := in.Create(repo, clock)
	if err == nil {
		desc := "returns an error when repo errored"
		t.Errorf("%s: Create(_, %v, %v, %v) = (_, nil), want (_, error)", desc, repo, clock, in)
	}
}

func TestNewsToRow(t *testing.T) {
	inID := 123
	in := &News{
		News: &news.News{
			Author:      "test-author",
			Title:       "test-title",
			Description: "test-description",
			URL:         "http://test-url",
			ImageURL:    "http://test-image-url",
			PublishedAt: time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	want := store.Row{
		inID,
		in.Author,
		in.Title,
		in.Description,
		in.URL,
		in.ImageURL,
		in.PublishedAt,
	}

	got := newsToRow(inID, in)
	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a store.Row given a news"
		t.Errorf("%s: newsToRow(%d, %v): Diff (-got +want)\n%s", desc, inID, in, diff)
	}
}

func TestSrcToRow(t *testing.T) {
	inID := 123
	in := &Source{
		Source: &news.Source{
			ID:   "test-id-456",
			Name: "test-source",
		},
	}

	want := store.Row{
		inID,
		in.ID,
		in.Name,
	}

	got := srcToRow(inID, in)
	if diff := pretty.Compare(got, want); diff != "" {
		desc := "returns a store.Row given a source"
		t.Errorf("%s: srcToRow(%d, %v): Diff (-got +want)\n%s", desc, inID, in, diff)
	}
}
