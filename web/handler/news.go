package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/riacataquian/news/api/news"
	"github.com/riacataquian/news/internal/auth"
	"github.com/riacataquian/news/internal/newsclient"
	"github.com/riacataquian/news/internal/newsclient/headlines"
	"github.com/riacataquian/news/internal/newsclient/list"

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
// TODO: Inject page and pageSize values to client params.
// If none exists, use default values in client.
func List(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	// Requests to external services should have timeouts.
	reqCtx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	client = newsclient.NewFromContext(reqCtx, listEndpoint)
	params := new(list.Params)
	err := schema.NewDecoder().Decode(params, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	res, err := fetch(params)
	if err != nil {
		return nil, err
	}

	// TODO: Missing page details.
	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.RequestURI,
		Count:      len(res.Articles),
		TotalCount: res.TotalResults,
		Data:       res.Articles,
	}, nil
}

// TopHeadlines is the HTTP handler for news requests to newsapi's top-headlines endpoint.
//
// Official docs: https://newsapi.org/docs/endpoints/top-headlines.
// TODO: Inject page and pageSize values to client params.
// If none exists, use default values in client.
func TopHeadlines(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	r.ParseForm()

	// Requests to external services should have timeouts.
	reqCtx, cancel := context.WithTimeout(ctx, defaultDuration)
	defer cancel()

	client = newsclient.NewFromContext(reqCtx, headlinesEndpoint)
	params := new(headlines.Params)
	err := schema.NewDecoder().Decode(params, r.Form)
	if err != nil {
		return nil, fmt.Errorf("error decoding params: %v", err)
	}

	res, err := fetch(params)
	if err != nil {
		return nil, err
	}

	// TODO: Missing page details.
	return &SuccessResponse{
		Code:       http.StatusOK,
		RequestURL: r.RequestURI,
		Count:      len(res.Articles),
		TotalCount: res.TotalResults,
		Data:       res.Articles,
	}, nil
}

// fetch performs the request to the client given params.
func fetch(params newsclient.Params) (*news.Response, error) {
	authKey, err := auth.LookupAndSetAuth()
	if err != nil {
		return nil, err
	}

	news, err := client.Get(authKey, params)
	if err != nil {
		return nil, err
	}

	// TODO: Inject query params from current request.
	//
	// qp := r.URL.Query()
	// ps := qp.Get("pageSize")
	// if ps != "" {
	// 	count, err := strconv.Atoi(ps)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid page size value: %v", ps)
	// 	}
	// 	resp.Count = count
	// }

	// p := qp.Get("page")
	// if p != "" {
	// 	page, err := strconv.Atoi(p)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("invalid page value: %v", p)
	// 	}
	// 	resp.Page = page
	// }

	return news, nil
}
