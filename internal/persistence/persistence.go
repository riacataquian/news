// Package persistence handles the processing for news list and the
// the interaction with the data repository.
package persistence

import (
	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/clock"
	"github.com/riacataquian/news/internal/store"
)

// News describes a news object.
type News struct {
	*news.News
}

// Source describes a news source.
type Source struct {
	*news.Source
}

// ScanRow scans a newsapi's news object into News.
func ScanRow(row *news.News) *News {
	return &News{
		News: &news.News{
			Source:      row.Source,
			Author:      row.Author,
			Title:       row.Title,
			Description: row.Description,
			URL:         row.URL,
			ImageURL:    row.ImageURL,
			PublishedAt: row.PublishedAt,
		},
	}
}

// Create persists rows to the supplied data repository.
// It makes use of timer to retrieve time.Now().Nanosecond() which is used as a resource ID.
func (row *News) Create(repo store.Store, timer clock.Time) error {
	nsecid := timer.Now().Nanosecond()

	var srow store.Row
	nrow := newsToRow(nsecid, row)
	if row.Source != nil {
		src := &Source{
			Source: &news.Source{
				ID:   row.Source.ID,
				Name: row.Source.Name,
			},
		}
		srow = srcToRow(nsecid, src)
	}

	// NOTE: Insertion is relative to the column declaration, order matters.

	nc := []string{"app_id", "author", "title", "description", "url", "image_url", "published_at"}
	rows := []store.Row{nrow}
	if err := repo.Create("news", nc, rows...); err != nil {
		return err
	}

	if row.Source != nil {
		sc := []string{"news_id", "id", "name"}
		rows := []store.Row{srow}
		if err := repo.Create("source", sc, rows...); err != nil {
			return err
		}
	}

	return nil
}

func newsToRow(nsecid int, n *News) (row store.Row) {
	row = append(row, nsecid)
	row = append(row, n.Author)
	row = append(row, n.Title)
	row = append(row, n.Description)
	row = append(row, n.URL)
	row = append(row, n.ImageURL)
	row = append(row, n.PublishedAt)
	return
}

func srcToRow(nsecid int, s *Source) (row store.Row) {
	row = append(row, nsecid)
	row = append(row, s.ID)
	row = append(row, s.Name)
	return
}
