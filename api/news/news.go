// Package news contains response mapping from https://newsapi.org.
package news

import "time"

// Response describes a succesful response from news API.
type Response struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []News `json:"articles"`
}

// News describes a news object.
type News struct {
	Source      `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

// Source describes a source object.
// It is the news object's source.
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse describes a failing response from news API.
//
// It satisfies the error interface.
type ErrorResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error is ErrorResponse's error interface implementation.
func (e ErrorResponse) Error() string {
	return e.Message
}
