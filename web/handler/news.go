package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/everything"
	"github.com/riacataquian/news/internal/newsclient/headlines"

	"github.com/gorilla/schema"
)

// This file contains handlers for news endpoint.

var client newsclient.Client

// List is the HTTP handler for news requests.
// See https://newsapi.org/docs/endpoints/everything for the official documentation.
func List(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = headlines.Client{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: everything.Endpoint,
		},
	}

	dst := new(headlines.Params)
	err := schema.NewDecoder().Decode(dst, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	return fetch(ctx, r, client, dst)
}

// TopHeadlines is the HTTP handler for top-headlines news requests.
// See https://newsapi.org/docs/endpoints/top-headlines for the official documentation.
func TopHeadlines(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = headlines.Client{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: headlines.Endpoint,
		},
	}

	dst := new(headlines.Params)
	err := schema.NewDecoder().Decode(dst, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	return fetch(ctx, r, client, dst)
}

// fetch performs the request to the client given params.
func fetch(ctx context.Context, r *http.Request, client newsclient.Client, params newsclient.Params) (*SuccessResponse, error) {
	news, err := client.Get(ctx, r, params)
	if err != nil {
		return nil, err
	}

	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.URL.String(),
		Data:       news,
	}, nil
}
