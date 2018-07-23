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

// News ...
func News(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = newsclient.NewsClient{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: newsclient.APIBaseURL + newsclient.TopHeadlinesPathPrefix,
		},
		ContextOrigin: ctx,
		RequestOrigin: r,
	}

	dst := new(newsclient.TopHeadlinesParams)
	err := schema.NewDecoder().Decode(dst, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	return fetchNews(client, dst)
}

// fetchNews ...
func fetchNews(client newsclient.Client, dst newsclient.Params) (*SuccessResponse, error) {
	news, err := client.GetTopHeadlines(dst)
	if err != nil {
		return nil, err
	}

	r := client.GetRequestOrigin()
	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.URL.String(),
		Data:       news,
	}, nil
}
