package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/riacataquian/news/internal/newsclient"

	"github.com/gorilla/schema"
)

// This file contains handlers for news endpoint.

var client newsclient.Client

// News is the HTTP handler for news requests.
func News(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = newsclient.NewsClient{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: newsclient.APIBaseURL + newsclient.TopHeadlinesPathPrefix,
		},
	}

	dst := new(newsclient.TopHeadlinesParams)
	err := schema.NewDecoder().Decode(dst, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	return fetchNews(ctx, r, client, dst)
}

// fetchNews ...
func fetchNews(ctx context.Context, r *http.Request, client newsclient.Client, params newsclient.Params) (*SuccessResponse, error) {
	news, err := client.TopHeadlines(ctx, r, params)
	if err != nil {
		return nil, err
	}

	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.URL.String(),
		Data:       news,
	}, nil
}
