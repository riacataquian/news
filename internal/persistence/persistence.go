// Package persistence handles the processing for news list and the
// the interaction with the data repository.
package persistence

import (
	"context"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/clock"
	"github.com/riacataquian/news/internal/store"
)

// Create persists rows to the supplied data repository.
// It makes use of timer to retrieve time.Now().Nanosecond() which is used as a resource ID.
func Create(_ context.Context, repo store.Store, timer clock.Time, rows []*news.News) error {
	var nrs []store.Row
	var srs []store.Row
	for _, news := range rows {
		nsecid := timer.Now().Nanosecond()
		nrs = append(nrs, newsToRow(nsecid, news))

		if news.Source != nil {
			srs = append(srs, srcToRow(nsecid, news.Source))
		}
	}

	// NOTE: Insertion is relative to the column declaration, order matters.

	nc := []string{"app_id", "author", "title", "description", "url", "image_url", "published_at"}
	if err := repo.Create("news", nc, nrs...); err != nil {
		return err
	}

	if len(srs) > 0 {
		sc := []string{"news_id", "id", "name"}
		if err := repo.Create("source", sc, srs...); err != nil {
			return err
		}
	}

	return nil
}

func newsToRow(nsecid int, n *news.News) (row store.Row) {
	row = append(row, nsecid)
	row = append(row, n.Author)
	row = append(row, n.Title)
	row = append(row, n.Description)
	row = append(row, n.URL)
	row = append(row, n.ImageURL)
	row = append(row, n.PublishedAt)
	return
}

func srcToRow(nsecid int, s *news.Source) (row store.Row) {
	row = append(row, nsecid)
	row = append(row, s.ID)
	row = append(row, s.Name)
	return
}
