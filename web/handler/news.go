package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/auth"
	"github.com/riacataquian/news/internal/httperror"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/headlines"
	"github.com/riacataquian/news/internal/newsclient/list"
	"github.com/riacataquian/news/internal/store"

	"github.com/gorilla/schema"
)

// This file contains handlers for news endpoint.

var (
	client newsclient.HTTPClient

	defaultDuration = 5 * time.Second

	listEndpoint      = list.ServiceEndpoint
	headlinesEndpoint = headlines.ServiceEndpoint
)

// List is the HTTP handler for news requests to newsapi's everything endpoint.
//
// Official docs: https://newsapi.org/docs/endpoints/everything.
func List(ctx context.Context, _ store.Store, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = newsclient.New(listEndpoint)
	params := new(list.Params)
	err := schema.NewDecoder().Decode(params, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	if params.Page == 0 {
		params.Page = 1
	}

	// Requests to external services should have timeouts.
	reqCtx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	res, err := fetch(reqCtx, params)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			RequestURL: r.RequestURI,
			DocsURL:    listEndpoint.DocsURL,
		}
	}

	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.RequestURI,
		Count:      len(res.Articles),
		Page:       params.Page,
		TotalCount: res.TotalResults,
		Data:       res.Articles,
	}, nil
}

// TopHeadlines is the HTTP handler for news requests to newsapi's top-headlines endpoint.
//
// Official docs: https://newsapi.org/docs/endpoints/top-headlines.
func TopHeadlines(ctx context.Context, _ store.Store, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	client = newsclient.New(headlinesEndpoint)
	params := new(headlines.Params)
	err := schema.NewDecoder().Decode(params, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	if params.Page == 0 {
		params.Page = 1
	}

	// Requests to external services should have timeouts.
	reqCtx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	res, err := fetch(reqCtx, params)
	if err != nil {
		return nil, &httperror.HTTPError{
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			RequestURL: r.RequestURI,
			DocsURL:    headlinesEndpoint.DocsURL,
		}
	}

	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.RequestURI,
		Count:      len(res.Articles),
		Page:       params.Page,
		TotalCount: res.TotalResults,
		Data:       res.Articles,
	}, nil
}

// fetch performs the request to the client given params.
func fetch(ctx context.Context, params newsclient.Params) (*news.Response, error) {
	authKey, err := auth.LookupAPIAuthKey()
	if err != nil {
		return nil, err
	}

	return client.Get(ctx, authKey, params)
}
