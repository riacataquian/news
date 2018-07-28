package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/headlines"
	"github.com/riacataquian/news/internal/newsclient/list"

	"github.com/gorilla/schema"
)

// This file contains handlers for news endpoint.

var client newsclient.Client

// List is the HTTP handler for news requests.
// See https://newsapi.org/docs/endpoints/everything for the official documentation.
func List(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = list.Client{
		ServiceEndpoint: newsclient.ServiceEndpoint{
			URL: list.Endpoint,
		},
	}

	dst := new(list.Params)
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
	resp := &SuccessResponse{Data: news}

	qp := r.URL.Query()
	ps := qp.Get("pageSize")
	if ps != "" {
		count, err := strconv.Atoi(ps)
		if err != nil {
			return nil, fmt.Errorf("invalid page size value: %v", ps)
		}
		resp.Count = count
	}

	p := qp.Get("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid page value: %v", p)
		}
		resp.Page = page
	}

	resp.Code = http.StatusOK
	resp.RequestURL = r.URL.String()
	return resp, nil
}
