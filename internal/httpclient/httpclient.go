package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/riacataquian/news/api/news"
)

// Dispatcher ...
type Dispatcher interface {
	DispatchRequest(*http.Request) (*news.Response, error)
}

// Client ...
type Client struct{}

// NewClient ...
func NewClient() Dispatcher {
	return &Client{}
}

// DispatchRequest dispatches given r http.Request.
//
// It encodes and return a news.ErrorResponse when an error is encountered.
// Returns news.Response otherwise for successful requests.
func (c *Client) DispatchRequest(r *http.Request) (*news.Response, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var res news.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}
		return nil, &res
	}

	var res news.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &res, nil
}
