// Package news ...
// TODO: Encode me to proto. For now, JSON will suffice.
package news

import "time"

// Response ...
type Response struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []News `json:"articles"`
}

// News ...
type News struct {
	Source      `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

// Source ...
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse ...
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
