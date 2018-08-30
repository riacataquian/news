// Package news contains response mapping and definitions for https://newsapi.org endpoint responses.
// Types here are meant to be consumed by a client.
//
// See Response Object > https://newsapi.org/docs/endpoints/everything.
package news

import "time"

// Response describes a successful response from newsapi.
type Response struct {
	Status string `json:"status"`
	// TotalResults are the total count of queryable results.
	// Use page parameter to page through the results.
	TotalResults int     `json:"totalResults"`
	Articles     []*News `json:"articles"`
}

// News describes a news object from newsapi.
type News struct {
	*Source     `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"` // in UTC format.
}

// Source describes a news source from newsapi.
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse describes a failing response from newsapi.
//
// It satisfies the error interface.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error is ErrorResponse's error interface implementation.
func (e *ErrorResponse) Error() string {
	return e.Message
}
